help:
	@echo "Available commands:"
	@echo " make docker-build                 Build the services"
	@echo " make docker-up                    Run the services"
	@echo " make docker-down                  Drop the services"
	@echo " make migrate-up                   Apply all database migrations"
	@echo " make migrate-down                 Rollback all database migrations"

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	docker compose run --rm auth-migrator -up

migrate-down:
	docker compose run --rm auth-migrator -down 0
