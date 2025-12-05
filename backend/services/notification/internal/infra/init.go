package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/pasarly/backend/services/notification/configs"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/subscriber"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/tracer"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type Infra struct {
	config   *configs.Config
	database *pgxpool.Pool
	logger   *zap.Logger
	mailer   *gomail.Dialer
	tracer   *tracer.Tracer

	acs *kafka.Reader
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		return nil, err
	}

	m := mailer.Init(&cfg.Mailer, l)

	t, err := tracer.Init(cfg.App.Name, cfg.Tracer.Endpoint, l)
	if err != nil {
		return nil, err
	}

	// Subscribers
	acs := subscriber.Init(&cfg.Broker, constants.EventTopicAuthCreated, l)

	return &Infra{
		config:   cfg,
		database: db,
		logger:   l,
		mailer:   m,
		tracer:   t,
		acs:      acs,
	}, nil
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) Mailer() *gomail.Dialer {
	return i.mailer
}

func (i *Infra) SubAuthCreated() *kafka.Reader {
	return i.acs
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.acs.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber (%s): %w", constants.EventTopicAuthCreated, err)
	}

	i.database.Close()
	i.tracer.Cleanup()
	return nil
}
