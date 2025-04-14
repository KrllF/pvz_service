package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/db"
	"github.com/KrllF/pvz_service/internal/handler/grpc"
	"github.com/KrllF/pvz_service/internal/handler/httph"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/repository/incache"
	"github.com/KrllF/pvz_service/internal/repository/postgres"
	"github.com/KrllF/pvz_service/internal/repository/postgreslog"
	"github.com/KrllF/pvz_service/internal/repository/txmanager"
	"github.com/KrllF/pvz_service/internal/server"
	"github.com/KrllF/pvz_service/internal/service/order/accept"
	audite "github.com/KrllF/pvz_service/internal/service/order/audit"
	"github.com/KrllF/pvz_service/internal/service/order/cache"
	"github.com/KrllF/pvz_service/internal/service/order/clusterkafka/consumer"
	"github.com/KrllF/pvz_service/internal/service/order/clusterkafka/producer"
	"github.com/KrllF/pvz_service/internal/service/order/cron/brocker"
	cronlog "github.com/KrllF/pvz_service/internal/service/order/cron/log"
	"github.com/KrllF/pvz_service/internal/service/order/returns"
	"github.com/KrllF/pvz_service/internal/service/order/show"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/KrllF/pvz_service/pkg/tracing"
	"go.uber.org/zap"
)

const capacity = 1000

type (
	// Server определяет методы для работы с обработчиком
	Server interface {
		Run() error
	}
	// Audit сервис
	Audit interface {
		Run(ctx context.Context)
	}
	// Cron interface
	Cron interface {
		Run(ctx context.Context) error
	}
	// Consumer interface
	Consumer interface {
		Run(ctx context.Context) error
	}
	// App структура - определяет App
	App struct {
		servHTTP    Server
		servGRPC    Server
		audit       Audit
		cron        Cron
		cronBrocker Cron
		consumer    Consumer
		traceClose  io.Closer
		closer      []Closer
	}
	// Closer определяет методы для закрытия
	Closer interface {
		Close()
	}
)

// NewApp создаёт новый экземляр AppClose()
// Принимает путь к JSON файлу - возвращает указатель на созданный App или ошибку
func NewApp(ctx context.Context) (*App, error) {
	conf, err := config.New("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config.New: %w", err)
	}
	_, closer := tracing.InitTracer("pvz-service")

	logg, err := logger.NewLogger(zap.InfoLevel)
	if err != nil {
		return nil, fmt.Errorf("logger.NewLogger: %w", err)
	}
	db, err := db.NewDB(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("db.NewDB: %w", err)
	}
	tx, err := txmanager.NewTxManager(db)
	if err != nil {
		return nil, fmt.Errorf("txmanager.NewTxManager: %w", err)
	}
	repo, err := postgres.NewRepository(tx, logg)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewRepository: %w", err)
	}
	repolog, err := postgreslog.NewRepository(tx, logg)
	if err != nil {
		return nil, fmt.Errorf("postgreslog.NewRepository: %w", err)
	}
	cche := incache.NewRepository[models.Order](int64(capacity))
	if _, err = cache.NewCacheService(ctx, repo, cche, true); err != nil {
		return nil, fmt.Errorf("cache.NewCacheService: %w", err)
	}

	ccheGet := incache.NewRepository[models.Order](int64(capacity))
	if _, err = cache.NewCacheService(ctx, repo, ccheGet, false); err != nil {
		return nil, fmt.Errorf("cache.NewCacheService: %w", err)
	}
	produc, err := producer.NewProducer("localhost:9092")
	if err != nil {
		return nil, fmt.Errorf("producer.NewProducer: %w", err)
	}
	crnBrocker, err := brocker.NewCron(repolog, produc, tx, conf.KafkaTopic, logg)
	if err != nil {
		return nil, fmt.Errorf("brocker.NewCron: %w", err)
	}
	crn, err := cronlog.NewCron(repo, ccheGet, logg)
	if err != nil {
		return nil, fmt.Errorf("cron.NewCron: %w", err)
	}
	consum, err := consumer.NewConsumer(conf)
	if err != nil {
		return nil, fmt.Errorf("consumer.NewConsumer: %w", err)
	}

	cls := make([]Closer, 0, 1)
	accServ := accept.NewService(repo, cche, tx, logg)
	retServ := returns.NewService(repo, cche, tx, logg)
	showServ, err := show.NewService(repo, ccheGet, logg)
	if err != nil {
		return nil, fmt.Errorf("show.NewService: %w", err)
	}
	auditServ := audite.NewAudite(repolog, tx, conf, logg)

	handHTTP := httph.NewHandler(ctx, accServ, retServ, showServ, auditServ)
	servHTTP := server.NewServer(conf, handHTTP.Init(conf))

	handGRPC := grpc.NewHandler(accServ, retServ, showServ, logg)
	servGRPC := server.NewServerGRPC(conf, handGRPC, auditServ)
	cls = append(cls, db)
	cls = append(cls, servHTTP)
	cls = append(cls, servGRPC)
	cls = append(cls, auditServ)
	cls = append(cls, produc)
	cls = append(cls, consum)

	return &App{
		servHTTP: servHTTP, servGRPC: servGRPC,
		closer: cls, traceClose: closer, audit: auditServ, cron: crn, cronBrocker: crnBrocker, consumer: consum,
	}, nil
}

// Run запустить
func (app *App) Run() error {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Printf("Cервер HTTP прослушивается...")
		if err := app.servHTTP.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	go func() {
		log.Printf("Cервер GRPC прослушивается...")
		if err := app.servGRPC.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	go func() {
		app.audit.Run(ctx)
	}()

	go func() {
		if err := app.cronBrocker.Run(ctx); err != nil {
			log.Fatalf("Ошибка при запуске крона: %v", err)
		}
	}()

	go func() {
		if err := app.cron.Run(ctx); err != nil {
			log.Fatalf("Ошибка при запуске крона: %v", err)
		}
	}()

	go func() {
		if err := app.consumer.Run(ctx); err != nil {
			log.Fatalf("ошибка при запуске consumer: %v", err)
		}
	}()

	sig := <-exit
	log.Printf("Получен сигнал завершения: %v", sig)
	cancel()

	app.traceClose.Close()
	for _, closer := range app.closer {
		closer.Close()
		log.Println("Успешно закрыто")
	}

	log.Println("Закрытие ресурсов прошло успешно")

	return nil
}
