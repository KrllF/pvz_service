//go:build integration

//nolint:all
package httptest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	myConf "github.com/KrllF/pvz_service/internal/config"
	dataB "github.com/KrllF/pvz_service/internal/db"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"

	"github.com/KrllF/pvz_service/internal/handler/httph"
	"github.com/KrllF/pvz_service/internal/repository/postgres"
	"github.com/KrllF/pvz_service/internal/server"
	"github.com/KrllF/pvz_service/internal/service/order/accept"
	"github.com/KrllF/pvz_service/internal/service/order/returns"
	"github.com/KrllF/pvz_service/internal/service/order/show"
)

var (
	DockerCli    *client.Client
	tdb          TDB
	TestServer   *server.Server
	Conffile     myConf.Config
	DBMigrations string
)

// TestMain выполняет глобальную инициализацию и очистку
func TestMain(m *testing.M) {
	if err := Setup(); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := TearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}

// Setup выполняет настройку перед запуском тестов
func Setup() error {
	ctx := context.Background()

	conffile, err := myConf.New("config.yaml")
	if err != nil {
		return fmt.Errorf("не удалось загрузить конфиг: %w", err)
	}
	Conffile = conffile

	var dockerErr error
	dockerCli, dockerErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if dockerErr != nil {
		return fmt.Errorf("не удалось создать докер клиента: %w", dockerErr)
	}
	DockerCli = dockerCli

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("ошибка при загрузке .env файла: %w", err)
	}

	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbMigrations := os.Getenv("MIGRATION_DIR")
	DBMigrations = dbMigrations
	imageName := os.Getenv("IMAGE_NAME")
	containerName := os.Getenv("CONTAINER_NAME")

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

	removeContainer(ctx, containerName)

	resp, err := dockerCli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("ошибка при создании контейнера: %v", err)
	}

	if err := dockerCli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("ошибка при запуске контейнера: %v", err)
	}

	if err = waitForDB(conffile.ConnectConfig.DSN, 30*time.Second); err != nil {
		return fmt.Errorf("не удалось подключится к базе: %w", err)
	}
	dat, err := sql.Open("postgres", conffile.ConnectConfig.DSN)
	if err != nil {
		return fmt.Errorf("не получилось подключиться к БД: %v", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("не удалось установить диалект: %w", err)
	}
	if err := goose.Up(dat, dbMigrations); err != nil {
		dat.Close()
		return fmt.Errorf("не получилось накатить миграции: %v", err)
	}
	dat.Close()

	db, dbErr := dataB.NewDB(context.Background(), conffile)
	if dbErr != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", dbErr)
	}
	tdb.DB = db
	repo, err := postgres.NewRepository(db)
	if err != nil {
		return fmt.Errorf("не получилось инициализировать репозиторий: %v", err)
	}

	accService := accept.NewService(repo)
	retService := returns.NewService(repo)
	showService := show.NewService(repo)

	handl := httph.NewHandler(ctx, accService, retService, showService)

	testServer := server.NewServer(conffile, handl.Init(conffile))
	TestServer = testServer
	go func() {
		if err := testServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ошибка при работе с сервером: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	return nil
}

// TearDown выполняет очистку после завершения тестов
func TearDown() error {
	ctx := context.Background()

	TestServer.Close()
	tdb.DB.Close()

	dat, err := sql.Open("postgres", Conffile.ConnectConfig.DSN)
	if err != nil {
		return fmt.Errorf("не получилось подключиться к БД: %v", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("не удалось установить диалект: %w", err)
	}
	if err := goose.Down(dat, DBMigrations); err != nil {
		dat.Close()
		return fmt.Errorf("не получилось накатить миграции: %v", err)
	}
	if err := goose.Down(dat, DBMigrations); err != nil {
		dat.Close()
		return fmt.Errorf("не получилось накатить миграции: %v", err)
	}
	dat.Close()

	if DockerCli != nil {
		containerName := os.Getenv("CONTAINER_NAME")
		removeContainer(ctx, containerName)
	}

	return nil
}

// removeContainer удаляет контейнер
func removeContainer(ctx context.Context, containerName string) {
	if DockerCli == nil {
		return
	}

	containers, err := DockerCli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Printf("Не удалось получить список контейнеров: %v", err)

		return
	}

	for _, contain := range containers {
		if len(contain.Names) > 0 && contain.Names[0] == "/"+containerName {
			if err := DockerCli.ContainerStop(ctx, contain.ID, container.StopOptions{}); err != nil {
				log.Printf("Не удалось остановить контейнер: %v", err)
			}
			if err := DockerCli.ContainerRemove(ctx, contain.ID, container.RemoveOptions{Force: true}); err != nil {
				log.Printf("Не удалось удалить контейнер: %v", err)
			}

			break
		}
	}
}

// waitForDB ожидание коннекта к базе
func waitForDB(dsn string, timeout time.Duration) error {
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
