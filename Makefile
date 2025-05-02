include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

lines := 1000
compose := docker compose
LOCAL_BIN:=$(CURDIR)/bin

install-debs:
	mkdir -p bin
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

init:
	cp env.example .env

lint:
	golangci-lint run --timeout=5m

create-migration:
	@read -p "Enter migration name: " name; \
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	$(LOCAL_BIN)/goose -dir=migrations create $${TIMESTAMP}_$${name} sql
