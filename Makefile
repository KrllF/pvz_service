include .env

APP_NAME = pvz
MAIN_FILE = cmd/app/main.go
LINT_FILE = .golangci.yaml
LOCAL_BIN = $(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN=$(PG_DSN)

build: lint
	go build -o $(APP_NAME) $(MAIN_FILE)
# Установка goose локально
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
		GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0


generate:
	make generate-order-api

generate-order-api:
	mkdir -p pkg/order_v1
	protoc --proto_path api/order_v1 \
	--experimental_allow_proto3_optional \
	--go_out=pkg/order_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/order_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/order_v1/order.proto

# Запуск приложения (со сборкой)
run: lint build
	./$(APP_NAME)

# Запуск приложения (без сборки)
run-only:
	./$(APP_NAME)

clear:
	rm -rf $(LOCAL_BIN)
	rm $(APP_NAME)

# ОБновляем зависимости
deps:
	go mod tidy

# Качаем зависимости
install: deps
	go mod download

# Запуск unit-tests
unit-tests:
	go test -v ./... -tags=unit

# Запуск integration-test
integration-tests:
	go test -v ./... -tags=integration

# Запуск suite-test
suite-tests:
	go test -v ./... -tags=suite

# линтер
lint:
	golangci-lint run -c $(LINT_FILE)

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
