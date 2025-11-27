package logger

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"go.uber.org/zap"
)

func Init(cfg *configs.App) (*zap.Logger, error) {
	var l *zap.Logger
	var err error

	if cfg.Env == "prod" {
		l, err = zap.NewProduction(zap.AddCaller())
	} else {
		l, err = zap.NewDevelopment(zap.AddCaller())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	l.Sugar().Infof("✅ [LOGGER] initialized (env=%s, level=%s)", cfg.Env, l.Level().String())
	return l, nil
}
