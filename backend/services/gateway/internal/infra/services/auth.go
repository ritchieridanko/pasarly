package services

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthService(cfg *configs.Service, l *zap.Logger) (apis.AuthServiceClient, error) {
	conn, err := grpc.NewClient(cfg.Auth.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}

	l.Sugar().Infof("✅ [AUTH-SERVICE] running on (host=%s, port=%d)", cfg.Auth.Host, cfg.Auth.Port)
	return apis.NewAuthServiceClient(conn), nil
}
