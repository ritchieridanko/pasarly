package services

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthService(cfg *configs.Service) (apis.AuthServiceClient, error) {
	conn, err := grpc.NewClient(cfg.Auth.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth service: %w", err)
	}

	return apis.NewAuthServiceClient(conn), nil
}
