# === Variables ===
BINARY_DIR := bin
APP_BIN := $(BINARY_DIR)/app
MIGRATOR_BIN := $(BINARY_DIR)/migrator
PROTOC := protoc
PROTO_DIR := api/proto/v1
PB_DIR := internal/interface/grpc/protobufs

help:
	@echo "Available commands:"
	@echo " make run-app                      Run the application locally"
	@echo " make build-app                    Build the application locally"
	@echo " make build-and-run-app            Build and run the application locally"
	@echo " make build-protobuf               Build the protobuf files"
	@echo " make drop-protobuf                Drop the built protobuf files"
	@echo " make setup-database               Create the database locally"
	@echo " make drop-database                Drop the database locally"
	@echo " make migrator-up                  Apply migrations locally"
	@echo " make migrator-down                Rollback N migrations locally"
	@echo " make migrator-down-all            Rollback all migrations locally"
	@echo " make build-migrator               Build the migrator binary locally"
	@echo " make build-and-run-migrator       Build and run the migrator locally"
	@echo " make docker-build                 Build the containers"
	@echo " make docker-up                    Run the containers"
	@echo " make docker-down                  Drop the containers"
	@echo " make docker-migrator-up           Run the migrator container up to apply migrations"
	@echo " make docker-migrator-down         Run the migrator container down to rollback all migrations"

# ---------- App Commands ----------
run-app:
	make build-protobuf
	go run cmd/app/main.go

build-app:
	make build-protobuf
	go build -o $(APP_BIN) cmd/app/main.go

build-and-run-app:
	make build-app
	./$(APP_BIN)

# ---------- Protobuf Commands ----------
build-protobuf:
	make drop-protobuf
	$(PROTOC) \
	--proto_path=$(PROTO_DIR) \
	--go_out=$(PB_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/*.proto

drop-protobuf:
	@find $(PB_DIR) -name "*.pb.go" -type f -delete

# ---------- Database Commands (local only) ----------
setup-database:
	psql -U postgres -h localhost -tc "SELECT 1 FROM pg_database WHERE datname = 'pasarly_auth_db'" | grep -q 1 || \
	psql -U postgres -h localhost -c "CREATE DATABASE pasarly_auth_db"

drop-database:
	psql -U postgres -h localhost -c "DROP DATABASE IF EXISTS pasarly_auth_db"

# ---------- Migration Commands (local only) ----------
migrator-up:
	go run cmd/migrator/main.go -up

migrator-down:
	go run cmd/migrator/main.go -down $(steps)

migrator-down-all:
	go run cmd/migrator/main.go -down 0

build-migrator:
	go build -o $(MIGRATOR_BIN) cmd/migrator/main.go

build-and-run-migrator:
	make build-migrator
	./$(MIGRATOR_BIN) -up

# ---------- Docker Commands ----------
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-migrator-up:
	docker compose run --rm auth-migrator -up

docker-migrator-down:
	docker compose run --rm auth-migrator -down 0
