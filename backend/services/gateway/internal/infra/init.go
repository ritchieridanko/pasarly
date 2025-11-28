package infra

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/services"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/tracer"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"go.uber.org/zap"
)

type Infra struct {
	config *configs.Config
	logger *zap.Logger
	tracer *tracer.Tracer
	as     apis.AuthServiceClient
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(&cfg.App)
	if err != nil {
		return nil, err
	}

	t, err := tracer.Init(cfg, l)
	if err != nil {
		return nil, err
	}

	// Services
	as, err := services.NewAuthService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}

	return &Infra{config: cfg, logger: l, tracer: t, as: as}, nil
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) AuthService() apis.AuthServiceClient {
	return i.as
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to flush logger: %w", err)
	}

	i.tracer.Cleanup()
	return nil
}
