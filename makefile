.PHONY: migrate build all
APP_NAME = tempgalias

all: migrate build

build:
	go build -o $(APP_NAME) ./src

run:
	go run ./src/main.go

migrate-create:
	go run ./migrations/cmd/main.go create $$name go

migrate:
	go run ./migrations/cmd/main.go up

migrate-down:
	go run ./migrations/cmd/main.go down