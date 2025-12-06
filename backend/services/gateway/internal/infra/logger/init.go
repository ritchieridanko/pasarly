package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Init(env string) (*zap.Logger, error) {
	var l *zap.Logger
	var err error

	if env == "prod" {
		l, err = zap.NewProduction(zap.AddCaller())
	} else {
		l, err = zap.NewDevelopment(zap.AddCaller())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	l.Sugar().Infof("✅ [LOGGER] initialized (env=%s, level=%s)", env, l.Level().String())
	return l, nil
}
