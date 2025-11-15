package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritchieridanko/pasarly/backend/services/auth-service/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/di"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/interface/grpc/server"
)

func main() {
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}

	infra, err := infra.Init(cfg)
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
	defer infra.Close()

	container := di.Init(cfg, infra)
	gs := container.GRPCServer()

	// Run the grpc server
	go func(gs *server.GRPCServer) {
		if err := gs.Start(); err != nil {
			log.Fatalln("FATAL ->", err.Error())
		}
	}(gs)

	// Handle app shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("[AUTH-SERVICE] shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer cancel()

	if err := gs.Shutdown(ctx); err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
}
