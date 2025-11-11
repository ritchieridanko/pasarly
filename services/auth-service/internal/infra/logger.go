package infra

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"go.uber.org/zap"
)

type Logger struct {
	baseLogger  *zap.Logger
	sugarLogger *zap.SugaredLogger
}

func NewLogger(cfg *configs.App) (*Logger, error) {
	var l *zap.Logger
	var err error

	if cfg.Env == "prod" {
		l, err = zap.NewProduction(zap.AddCaller())
	} else {
		l, err = zap.NewDevelopment(zap.AddCaller())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build a logger: %w", err)
	}

	s := l.Sugar()
	s.Infof("✅ [LOGGER] initialized (env=%s, level=%s)", cfg.Env, l.Level().String())
	return &Logger{baseLogger: l, sugarLogger: s}, nil
}

func (l *Logger) Base() *zap.Logger {
	return l.baseLogger
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugarLogger
}

func (l *Logger) Close() error {
	if err := l.baseLogger.Sync(); err != nil {
		return fmt.Errorf("failed to flush logger: %w", err)
	}
	return nil
}
