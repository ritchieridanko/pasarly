package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/di"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/interface/server"
)

func main() {
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}

	i, err := infra.Init(cfg)
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
	defer i.Close()

	container := di.Init(cfg, i)
	s := container.Server()

	// Run the server
	go func(s *server.Server) {
		if err := s.Start(); err != nil {
			log.Fatalln("FATAL ->", err.Error())
		}
	}(s)

	// Handle app shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("🛑 [%s] is shutting down...", cfg.App.Name)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
}
