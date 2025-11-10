# === Variables ===
BINARY_DIR := bin
MIGRATE_BIN := $(BINARY_DIR)/migrate

help:
	@echo "Available commands:"
	@echo " make setup-database               Create the database locally"
	@echo " make drop-database                Drop the database locally"
	@echo " make migrate-up                   Apply migrations locally"
	@echo " make migrate-down                 Rollback N migrations locally"
	@echo " make migrate-down-all             Rollback all migrations locally"
	@echo " make build-migrate                Build the migrate binary locally"
	@echo " make build-and-run-migrate        Build and run the migrate locally"
	@echo " make docker-build                 Build the containers"
	@echo " make docker-up                    Run the containers"
	@echo " make docker-down                  Drop the containers"
	@echo " make docker-migrate-up            Run the migrate container up to apply migrations"
	@echo " make docker-migrate-down          Run the migrate container down to rollback all migrations"

# ---------- Database Commands (local only) ----------
setup-database:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'pasarly_auth_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE pasarly_auth_db"

drop-database:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS pasarly_auth_db"

# ---------- Migration Commands (local only) ----------
migrate-up:
	go run cmd/migrate/main.go -up

migrate-down:
	go run cmd/migrate/main.go -down $(steps)

migrate-down-all:
	go run cmd/migrate/main.go -down 0

build-migrate:
	go build -o $(MIGRATE_BIN) cmd/migrate/main.go

build-and-run-migrate:
	make build-migrate
	./$(MIGRATE_BIN) -up

# ---------- Docker Commands ----------
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-migrate-up:
	docker compose run --rm auth-migrate -up

docker-migrate-down:
	docker compose run --rm auth-migrate -down 0
