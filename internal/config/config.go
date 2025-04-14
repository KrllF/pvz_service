package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type ConnectConfig struct {
	Login        string
	Password     string
	DSN          string
	KafkaTopic   string
	KafkaBrokers string
	HTTP_HOST    string
	HTTP_PORT    string
	GRPC_HOST    string
	GRPC_PORT    string

	PROMETHEUS_PORT string
}

type AppConfig struct {
	Read       time.Duration `yaml:"read"`
	Write      time.Duration `yaml:"write"`
	Idle       time.Duration `yaml:"idle"`
	ReadHeader time.Duration `yaml:"read_header"`
	Word       string        `yaml:"word"`
}

type Config struct {
	ConnectConfig
	AppConfig
}

func New(yamlPath string) (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Не удалось загрузить файл .env: %v", err)
		return Config{}, err
	}

	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Printf("Не удалось прочитать YAML-файл: %v", err)
		return Config{}, err
	}

	var appConfig AppConfig
	if err := yaml.Unmarshal(yamlFile, &appConfig); err != nil {
		log.Printf("Не удалось сделать Unmarshal в структуру: %v", err)
		return Config{}, err
	}

	return Config{
		ConnectConfig: ConnectConfig{
			Login:           getEnv("ADMIN_LOGIN"),
			Password:        getEnv("ADMIN_PASSWORD"),
			DSN:             getEnv("PG_DSN"),
			KafkaTopic:      getEnv("KAFKA_TOPIC"),
			KafkaBrokers:    getEnv("KAFKA_BROKERS"),
			HTTP_HOST:       getEnv("HTTP_HOST"),
			HTTP_PORT:       getEnv("HTTP_PORT"),
			GRPC_HOST:       getEnv("GRPC_HOST"),
			GRPC_PORT:       getEnv("GRPC_PORT"),
			PROMETHEUS_PORT: getEnv("PROMETHEUS_PORT"),
		},
		AppConfig: appConfig,
	}, nil
}
