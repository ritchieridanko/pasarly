package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ritchieridanko/pasarly/backend/services/user/configs"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/di"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/subscriber"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/interface/server"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/processors"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// Run the subscribers
	go func(ctx context.Context, s *subscriber.Subscriber, p processors.UserProcessor) {
		defer wg.Done()
		if err := s.Listen(ctx, p.OnAuthCreated); err != nil {
			log.Println("ERROR ->", err.Error())
		}
	}(ctx, container.SubAuthCreated(), container.UserProcessor())

	// Handle app shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("🛑 [%s] is shutting down...", cfg.App.Name)

	sdCtx, sdCancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer sdCancel()

	if err := s.Shutdown(sdCtx); err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}

	wg.Wait()
}
