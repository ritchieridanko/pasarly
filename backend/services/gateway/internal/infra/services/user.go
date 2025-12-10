package services

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserService(cfg *configs.Service, l *zap.Logger) (apis.UserServiceClient, error) {
	conn, err := grpc.NewClient(cfg.User.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user service: %w", err)
	}

	l.Sugar().Infof("✅ [USER-SERVICE] running on (host=%s, port=%d)", cfg.User.Host, cfg.User.Port)
	return apis.NewUserServiceClient(conn), nil
}
