package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ritchieridanko/pasarly/backend/services/user/configs"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"google.golang.org/grpc"
)

type Server struct {
	config *configs.Server
	server *grpc.Server
	logger *logger.Logger
}

func Init(cfg *configs.Server, l *logger.Logger, uh *handlers.UserHandler, ah *handlers.AddressHandler) *Server {
	s := grpc.NewServer()

	apis.RegisterUserServiceServer(s, uh)
	apis.RegisterUserAddressServiceServer(s, ah)

	return &Server{config: cfg, server: s, logger: l}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	if err := s.server.Serve(l); err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	s.logger.Sugar().Infof("✅ [SERVER] running on (host=%s, port=%d)", s.config.Host, s.config.Port)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return fmt.Errorf("failed to shutdown server: %w", ctx.Err())
	case <-stopped:
		return nil
	}
}
