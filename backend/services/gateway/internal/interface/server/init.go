package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
)

type Server struct {
	config *configs.Server
	server *http.Server
	logger *logger.Logger
}

func Init(cfg *configs.Server, h http.Handler, l *logger.Logger) *Server {
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      h,
		ReadTimeout:  cfg.Timeout.Read,
		WriteTimeout: cfg.Timeout.Write,
	}

	return &Server{config: cfg, server: s, logger: l}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	s.logger.Sugar().Infof("✅ [SERVER] running on (host=%s, port=%d)", s.config.Host, s.config.Port)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.server.Shutdown(ctx)
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Close()
		return fmt.Errorf("failed to shutdown server: %w", ctx.Err())
	case <-stopped:
		return nil
	}
}
