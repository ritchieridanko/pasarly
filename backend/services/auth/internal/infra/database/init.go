package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"go.uber.org/zap"
)

func Init(cfg *configs.Database, l *zap.Logger) (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	c.MaxConns = int32(cfg.MaxConns)
	c.MinConns = int32(cfg.MinConns)
	c.MaxConnLifetime = cfg.MaxConnLifetime
	c.MaxConnIdleTime = cfg.MaxConnIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	l.Sugar().Infof("✅ [DATABASE] initialized (host=%s, port=%d, name=%s)", cfg.Host, cfg.Port, cfg.Name)
	return pool, nil
}
