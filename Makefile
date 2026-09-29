APP=bin/server

ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: all build run start dev migrate-up migrate-down migrate-force migrate-reset migrate-create test fmt tidy install-tools
all: build

build:
	mkdir -p bin
	go build -o $(APP) ./cmd/server

run:
	go run ./cmd/server

start: build
	./$(APP)

dev:
	air -c .air.toml

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

migrate-reset:
	migrate -path migrations -database "$(DATABASE_URL)" drop -f

migrate-create:
	migrate create -ext sql -dir migrations $(name)

test:
	go test -v ./...

fmt:
	gofmt -w ./cmd ./internal

tidy:
	go mod tidy

install-tools:
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
