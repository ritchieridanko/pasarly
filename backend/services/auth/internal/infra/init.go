package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/tracer"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	cache    *redis.Client
	database *pgxpool.Pool
	logger   *zap.Logger
	tracer   *tracer.Tracer

	acp *kafka.Writer
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
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

	t, err := tracer.Init(cfg.App.Name, cfg.Tracer.Endpoint, l)
	if err != nil {
		return nil, err
	}

	// Publishers
	acp := publisher.Init(&cfg.Broker, constants.EventTopicAuthCreated, l)

	return &Infra{config: cfg, cache: c, database: db, logger: l, tracer: t, acp: acp}, nil
}

func (i *Infra) Cache() *redis.Client {
	return i.cache
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) PubAuthCreated() *kafka.Writer {
	return i.acp
}

func (i *Infra) Close() error {
	if err := i.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache: %w", err)
	}
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.acp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (%s): %w", constants.EventTopicAuthCreated, err)
	}

	i.database.Close()
	i.tracer.Cleanup()
	return nil
}
