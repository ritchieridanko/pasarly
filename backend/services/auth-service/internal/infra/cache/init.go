package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/configs"
	"go.uber.org/zap"
)

func Init(cfg *configs.Cache, l *zap.Logger) (*redis.Client, error) {
	if cfg.Pass == "" {
		l.Sugar().Warnln("⚠️ [CACHE] connecting without password...")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	c := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: cfg.Pass,
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("failed to ping cache: %w", err)
	}

	l.Sugar().Infof("✅ [CACHE] initialized (host=%s, port=%d)", cfg.Host, cfg.Port)
	return c, nil
}
