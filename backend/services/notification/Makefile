# === Variables ===
BINARY_DIR := bin
APP_BIN := $(BINARY_DIR)/app

help:
	@echo "Available commands:"
	@echo " make run-app                      Run the application"
	@echo " make build-app                    Build the application"
	@echo " make build-and-run-app            Build and run the application"

# ---------- App Commands ----------
run-app:
	go run cmd/app/main.go

build-app:
	go build -o $(APP_BIN) cmd/app/main.go

build-and-run-app:
	make build-app
	./$(APP_BIN)
