package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/auth-service/internal/interface/grpc/handlers"
	"github.com/ritchieridanko/pasarly/auth-service/internal/interface/grpc/protobufs/v1"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	config *configs.Server
	server *grpc.Server
	logger *logger.Logger
}

func NewGRPCServer(cfg *configs.Server, ah *handlers.AuthGRPCHandler, l *logger.Logger) *GRPCServer {
	s := grpc.NewServer()

	protobufs.RegisterAuthServiceServer(s, ah)

	return &GRPCServer{config: cfg, server: s, logger: l}
}

func (s *GRPCServer) Start() error {
	h := s.config.GRPC.Host
	p := s.config.GRPC.Port
	addr := fmt.Sprintf("%s:%d", h, p)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to build grpc listener: %w", err)
	}

	if err := s.server.Serve(l); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	s.logger.Sugar().Infof("✅ [GRPC SERVER] running on (host=%s, port=%d)", h, p)
	return nil
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}
