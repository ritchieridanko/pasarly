# === Variables ===
BINARY_DIR := bin
APP_BIN := $(BINARY_DIR)/app
MIGRATOR_BIN := $(BINARY_DIR)/migrator

help:
	@echo "Available commands:"
	@echo " make run-app                      Run the application"
	@echo " make build-app                    Build the application"
	@echo " make build-and-run-app            Build and run the application"
	@echo " make setup-database               Create the database"
	@echo " make drop-database                Drop the database"
	@echo " make migrate-up                   Apply all migrations"
	@echo " make migrate-down                 Rollback N migrations"
	@echo " make migrate-down-all             Rollback all migrations"
	@echo " make build-migrator               Build the migrator"
	@echo " make build-and-run-migrator       Build and run the migrator"

# ---------- App Commands ----------
run-app:
	go run cmd/app/main.go

build-app:
	go build -o $(APP_BIN) cmd/app/main.go

build-and-run-app:
	make build-app
	./$(APP_BIN)

# ---------- Database Commands ----------
setup-database:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'pasarly_auth_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE pasarly_auth_db"

drop-database:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS pasarly_auth_db"

# ---------- Migration Commands ----------
migrate-up:
	go run cmd/migrator/main.go -up

migrate-down:
	go run cmd/migrator/main.go -down $(steps)

migrate-down-all:
	go run cmd/migrator/main.go -down 0

build-migrator:
	go build -o $(MIGRATOR_BIN) cmd/migrator/main.go

build-and-run-migrator:
	make build-migrator
	./$(MIGRATOR_BIN) -up
