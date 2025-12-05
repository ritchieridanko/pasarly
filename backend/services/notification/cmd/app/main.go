package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ritchieridanko/pasarly/backend/services/notification/configs"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/di"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/subscriber"
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

	container, err := di.Init(cfg, i)
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
	defer container.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context, s *subscriber.Subscriber, h *handlers.AuthHandler) {
		defer wg.Done()
		if err := s.Listen(ctx, h.OnAuthCreated); err != nil {
			log.Println("ERROR ->", err.Error())
		}
	}(ctx, container.SubAuthCreated(), container.AuthHandler())

	// Handle app shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("🛑 [%s] is shutting down...", cfg.App.Name)
	cancel()
	wg.Wait()
}
