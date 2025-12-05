package di

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/notification/configs"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/channels"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/subscriber"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/repositories"
)

type Container struct {
	config   *configs.Config
	database *database.Database
	logger   *logger.Logger
	mailer   *mailer.Mailer
	acs      *subscriber.Subscriber
	ec       channels.EmailChannel
	er       repositories.EventRepository
	ah       *handlers.AuthHandler
}

func Init(cfg *configs.Config, i *infra.Infra) (*Container, error) {
	// Infra
	db := database.NewDatabase(i.Database())
	l := logger.NewLogger(i.Logger())
	m := mailer.NewMailer(i.Mailer())

	// Subscribers
	acs := subscriber.NewSubscriber(&cfg.Broker, i.SubAuthCreated(), l)

	// Channels
	ec, err := channels.NewEmailChannel(m, cfg.Client.BaseURL, cfg.Mailer.From)
	if err != nil {
		return nil, err
	}

	// Repositories
	er := repositories.NewEventRepository(db)

	// Handlers
	ah := handlers.NewAuthHandler(er, ec, cfg.Mailer.Timeout)

	return &Container{
		config:   cfg,
		database: db,
		logger:   l,
		mailer:   m,
		acs:      acs,
		ec:       ec,
		er:       er,
		ah:       ah,
	}, nil
}

func (c *Container) SubAuthCreated() *subscriber.Subscriber {
	return c.acs
}

func (c *Container) AuthHandler() *handlers.AuthHandler {
	return c.ah
}

func (c *Container) Close() error {
	if err := c.mailer.Close(); err != nil {
		return fmt.Errorf("failed to close mailer: %w", err)
	}
	return nil
}
