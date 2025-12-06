# === Variables ===
SHARED_DIR := shared
PROTO_DIR := $(SHARED_DIR)/protobufs
APIS_DIR := $(SHARED_DIR)/apis
EVENTS_DIR := $(SHARED_DIR)/events

help:
	@echo "Available commands:"
	@echo " make build-proto                  Build the proto files"
	@echo " make drop-proto                   Drop the built .pb.go files"
	@echo " make docker-build                 Build the services"
	@echo " make docker-up                    Run the services"
	@echo " make docker-down                  Drop the services"
	@echo " make migrate-up                   Apply all database migrations"
	@echo " make migrate-down                 Rollback all database migrations"

# ---------- Proto Commands ----------
build-proto:
	make drop-proto
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(APIS_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(APIS_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*/*_api.proto
	protoc \
	  --proto_path=$(PROTO_DIR) \
		--go_out=$(EVENTS_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(EVENTS_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*/*_event.proto

drop-proto:
	@find $(APIS_DIR) -name "*.pb.go" -type f -delete
	@find $(EVENTS_DIR) -name "*.pb.go" -type f -delete

# ---------- Docker Commands ----------
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down -v

migrate-up:
	docker compose run --rm auth-migrator -up
	docker compose run --rm user-migrator -up
	docker compose run --rm notification-migrator -up

migrate-down:
	docker compose run --rm auth-migrator -down 0
	docker compose run --rm user-migrator -down 0
	docker compose run --rm notification-migrator -down 0
