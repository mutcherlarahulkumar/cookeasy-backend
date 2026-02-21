APP_NAME := tweets-service

DB_USER := postgres
DB_PASSWORD := secret
DB_HOST := localhost
DB_PORT := 5432
DB_NAME := cookeasy

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

MIGRATIONS_PATH := db/migrations

up:
	docker compose up -d

down:
	docker compose down

down-v:
	docker compose down -v

restart: down up

logs:
	docker compose logs -f

ps:
	docker compose ps

run:
	go run cmd/main.go

build:
	go build -o bin/$(APP_NAME) cmd/main.go

test:
	go test ./...

tidy:
	go mod tidy

db-shell:
	docker exec -it db-go-tweets psql -U $(DB_USER) -d $(DB_NAME)

db-reset: down-v up


# ==============================
# Migrations (golang-migrate)
# Requires: brew install golang-migrate OR go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# ==============================

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-check:
	migrate -path db/migrations -database "$(DB_URL)" version



# ==============================
# Full Dev Helpers
# ==============================

dev: up migrate-up run

fresh: down-v up migrate-up run
