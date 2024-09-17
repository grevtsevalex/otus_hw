BIN := "./bin/calendar"
BIN_SC := "./bin/calendar_scheduler"
DOCKER_IMG="calendar:develop"
DOCKER_TEST_IMG="test_container:develop"
DOCKER_SCH_IMG="scheduler:develop"
DOCKER_SEN_IMG="sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# запустить все сервисы
up: down
	docker compose up --force-recreate

# остановить все сервисы
down:
	docker compose down 

# -- local -- 

# локальная сборка календаря 
local-calendar-build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
# локальный запуск календаря
local-calendar-run: local-calendar-build
	$(BIN) -config ./configs/config.toml

# локальная сборка планировщика 
local-scheduler-build: rabbit-up
	go build -v -o $(BIN_SC) -ldflags "$(LDFLAGS)" ./cmd/scheduler

# локальный запуск планировщика
local-scheduler-run: scheduler-build
	$(BIN_SC) -config ./configs/scheduler_config.toml

# локальная сборка отправщика 
local-sender-build:
	go build -v -o $(BIN_SC) -ldflags "$(LDFLAGS)" ./cmd/sender

# локальный запуск отправщика
local-sender-run: sender-build
	$(BIN_SC) -config ./configs/sender_config.toml

# локальный запуск рэббита в докере
rabbit-up:
	docker run --rm --hostname my-rabbit --name calendar-rabbit rabbitmq:3
	docker run -d --name rb1 -p 15672:15672 -p 5672:5672 rabbitmq:3-management

# -- -- --


# -- docker --

build: 
	make build-calendar-img 
	make build-scheduler-img 
	make build-sender-img

# сборка образа календаря
build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

# запуск контейнера календаря
run-calendar-img: build-calendar-img
	docker run $(DOCKER_IMG)

# сборка образа планировщика
build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SCH_IMG) \
		-f build/scheduler.Dockerfile .

# запуск контейнера планировщика
run-scheduler-img: build-scheduler-img
	docker run $(DOCKER_SCH_IMG)

# сборка образа отправщика
build-sender-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SEN_IMG) \
		-f build/sender.Dockerfile .

# запуск контейнера планировщика
run-sender-img: build-sender-img
	docker run $(DOCKER_SEN_IMG)


# сборка образа приложения тестирования
build-test-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_TEST_IMG) \
		-f tests/Dockerfile .

# запуск контейнера приложения тестирования
run-test-img: build-test-img
	docker run $(DOCKER_TEST_IMG)

# -- -- --


version: build
	$(BIN) version

# запустить тесты
test:
	go test -race ./internal/...

# скачать зависимости линтеров
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

# запустить линтеры
lint: install-lint-deps
	golangci-lint run ./...

integration-tests:
	docker compose -f ./docker-compose-test.yaml down
	docker compose -f ./docker-compose-test.yaml run test-container

.PHONY: build run build-img run-img version test lint

# выполнить миграцию
migrate:
	$(HOME)/go/bin/goose -dir ./migrations postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=disable password=$(DB_PASS)" up

# откатить миграцию
migrate-down:
	$(HOME)/go/bin/goose -dir ./migrations postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=disable password=$(DB_PASS)" down

# проверить статус миграций
migrate-status:
	$(HOME)/go/bin/goose -dir ./migrations postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=disable password=$(DB_PASS)" status

# сгенерить прото-файлы
generate:
	protoc api/eventService.proto --go_out=./internal/server/grpc/pb --go-grpc_out=./internal/server/grpc/pb
