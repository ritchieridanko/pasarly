package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/cache"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/publisher"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Infra struct {
	logger   *zap.Logger
	cache    *redis.Client
	database *pgxpool.Pool

	// Publishers
	acp *kafka.Writer
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(&cfg.App)
	if err != nil {
		return nil, err
	}

	c, err := cache.Init(&cfg.Cache, l)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		return nil, err
	}

	acp, err := publisher.Init(&cfg.Broker, "auth.created", l)
	if err != nil {
		return nil, err
	}

	return &Infra{logger: l, cache: c, database: db, acp: acp}, nil
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) Cache() *redis.Client {
	return i.cache
}

func (i *Infra) DB() *pgxpool.Pool {
	return i.database
}

func (i *Infra) PubAuthCreated() *kafka.Writer {
	return i.acp
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to flush logger: %w", err)
	}
	if err := i.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache: %w", err)
	}
	if err := i.acp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (auth.created): %w", err)
	}

	i.database.Close()
	return nil
}
