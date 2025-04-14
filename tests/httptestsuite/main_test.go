//go:build suite

package httptestsuite

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	myConf "github.com/KrllF/pvz_service/internal/config"
	dataB "github.com/KrllF/pvz_service/internal/db"
	"github.com/KrllF/pvz_service/internal/handler/httph"
	"github.com/KrllF/pvz_service/internal/repository/postgres"
	"github.com/KrllF/pvz_service/internal/server"
	"github.com/KrllF/pvz_service/internal/service/order/accept"
	"github.com/KrllF/pvz_service/internal/service/order/returns"
	"github.com/KrllF/pvz_service/internal/service/order/show"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"
)

var dbURI, dbName string

func init() {
	dbURI = os.Getenv("TEST_DB_URI")
	dbName = os.Getenv("TEST_DB_NAME")
}

type APITestSuite struct {
	suite.Suite

	db     DB
	server Server
	conf   myConf.Config

	dockerCLI    *client.Client
	migrDir      string
	contanerName string
	pool         *pgxpool.Pool
}

func NewAPITestSuite() *APITestSuite {
	return &APITestSuite{}
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, NewAPITestSuite())
}

func (s *APITestSuite) SetupSuite() {
	ctx := context.Background()
	conffile, err := myConf.New("config.yaml")
	if err != nil {
		s.FailNow("не получилось загрузить конфиг: %v", err)
	}
	s.conf = conffile

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		s.FailNow("не получилось создать докер клиента: %v", err)
	}
	s.dockerCLI = dockerCli

	if err := godotenv.Load(); err != nil {
		s.FailNow("не получилось загрузить .env: %v", err)
	}

	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbMigrations := os.Getenv("MIGRATION_DIR")
	s.migrDir = dbMigrations
	imageName := os.Getenv("IMAGE_NAME")
	containerName := os.Getenv("CONTAINER_NAME")
	s.contanerName = containerName

	env := []string{
		"POSTGRES_DB=" + dbName,
		"POSTGRES_USER=" + dbUser,
		"POSTGRES_PASSWORD=" + dbPassword,
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": {{HostIP: dbHost, HostPort: dbPort}},
		},
	}

	config := &container.Config{
		Image: imageName,
		Env:   env,
	}

	s.removeContainer(ctx, containerName)

	resp, err := dockerCli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		s.FailNow("не получилось создать контейнер: %v", err)
	}

	if err := dockerCli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		s.FailNow("не получилось запустить контейнер: %v", err)
	}

	if err = s.waitForDB(conffile.ConnectConfig.DSN, 30*time.Second); err != nil {
		s.FailNow("не получилось подключится к бд: %v", err)
	}
	log.Printf("база подключена")

	dat, err := sql.Open("postgres", conffile.ConnectConfig.DSN)
	if err != nil {
		s.FailNow("не получилось подключится к бд: %v", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		s.FailNow("не удалось установить диалект: %v", err)
	}
	if err := goose.Up(dat, dbMigrations); err != nil {
		dat.Close()
		s.FailNow("не получилось накатить миграции: %v", err)
	}
	dat.Close()

	db, err := dataB.NewDB(context.Background(), conffile)
	if err != nil {
		s.FailNow("не получилось подключится к бд: %v", err)
	}
	s.db = db
	s.pool = db.GetPool()
	repo, err := postgres.NewRepository(db)
	if err != nil {
		s.FailNow("не получилось инициализировать репозиторий: %v", err)
	}

	accService := accept.NewService(repo)
	retService := returns.NewService(repo)
	showService := show.NewService(repo)

	handl := httph.NewHandler(ctx, accService, retService, showService)
	testServer := server.NewServer(conffile, handl.Init(conffile))
	s.server = testServer
	go func() {
		if err := s.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ошибка при работе с сервером: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)
}

// SetupSubTest добавление и очищение таблицы перед тестом
func (s *APITestSuite) SetupSubTest() {
	s.truncateTable(context.Background())
	s.addUser(context.Background())
	s.addOrder(context.Background())
}

// TearDownSubTest очистка таблицы после теста
func (s *APITestSuite) TearDownSubTest() {
	s.truncateTable(context.Background())
}

func (s *APITestSuite) TearDownSuite() {
	ctx := context.Background()

	s.server.Close()
	s.db.Close()
	if s.pool != nil {
		s.pool.Close()
	}
	dat, err := sql.Open("postgres", s.conf.ConnectConfig.DSN)
	if err != nil {
		s.FailNow("не получилось подключится к бд: %v", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		s.FailNow("не удалось установить диалект: %v", err)
	}
	if err := goose.Down(dat, s.migrDir); err != nil {
		dat.Close()
		s.FailNow("не получилось откатить миграции: %v", err)
	}
	if err := goose.Down(dat, s.migrDir); err != nil {
		dat.Close()
		s.FailNow("не получилось откатить миграции: %v", err)
	}

	dat.Close()

	if s.dockerCLI != nil {
		s.removeContainer(ctx, s.contanerName)
	}
}

// removeContainer удаляет контейнер
func (s *APITestSuite) removeContainer(ctx context.Context, containerName string) {
	if s.dockerCLI == nil {
		return
	}

	containers, err := s.dockerCLI.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Printf("Не удалось получить список контейнеров: %v", err)

		return
	}

	for _, contain := range containers {
		if len(contain.Names) > 0 && contain.Names[0] == "/"+containerName {
			if err := s.dockerCLI.ContainerStop(ctx, contain.ID, container.StopOptions{}); err != nil {
				log.Printf("Не удалось остановить контейнер: %v", err)
			}
			if err := s.dockerCLI.ContainerRemove(ctx, contain.ID, container.RemoveOptions{Force: true}); err != nil {
				log.Printf("Не удалось удалить контейнер: %v", err)
			}

			break
		}
	}
}

// waitForDB ожидание коннекта к базе
func (s *APITestSuite) waitForDB(dsn string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("время ожидания подключения к базе вышло")
		case <-ticker.C:
			db, err := sql.Open("postgres", dsn)
			if err == nil {
				if err := db.Ping(); err == nil {
					db.Close()

					return nil
				}
			}
			db.Close()
		}
	}
}

// addUser добавить пользователя перед тестом
func (s *APITestSuite) addUser(ctx context.Context) {
	_, err := s.pool.Exec(ctx, "INSERT INTO users(user_id) VALUES($1)", 1)
	if err != nil {
		panic(err)
	}
}

// addOrder добавить заказ перед тестом
func (s *APITestSuite) addOrder(ctx context.Context) {
	queryAdd := `INSERT INTO orders(order_id,user_id,
	order_weight, order_price,
	pack_id, extra_pack_id,
	status_id, shelf_life,
	two_days_of_life, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8, $9, $10, $11);`

	twoDaysOfLife := pgtype.Timestamptz{
		Time:   time.Time{},
		Status: pgtype.Null,
	}

	_, err := s.pool.Exec(ctx, queryAdd, 1, 1, 100, 100,
		1, 1, 1, "2027-10-10T15:15:15Z", twoDaysOfLife, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}
	_, err = s.pool.Exec(ctx, queryAdd, 100, 1, 100, 100,
		1, 1, 3, "2023-10-10T15:15:15Z", twoDaysOfLife, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}
}

// truncateTable очистить таблицу
func (s *APITestSuite) truncateTable(ctx context.Context) {
	q := "TRUNCATE orders, users"
	if _, err := s.db.Exec(ctx, q); err != nil {
		panic(err)
	}
}
