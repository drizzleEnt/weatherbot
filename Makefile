lines := 1000
compose := docker compose
LOCAL_BIN:=$(CURDIR)/bin


init:
	cp env.example .env

lint:
	golangci-lint run --timeout=5m

create-migration:
	@read -p "Enter migration name: " name; \
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	$(LOCAL_BIN)/goose -dir=migrations create $${TIMESTAMP}_$${name} sql
