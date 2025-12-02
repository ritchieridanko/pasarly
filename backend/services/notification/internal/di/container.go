package di

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/notification/configs"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/channels"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/subscriber"
)

type Container struct {
	config *configs.Config
	logger *logger.Logger
	mailer *mailer.Mailer
	acs    *subscriber.Subscriber
	ec     channels.EmailChannel
	ah     *handlers.AuthHandler
}

func Init(cfg *configs.Config, i *infra.Infra) (*Container, error) {
	// Infra
	l := logger.NewLogger(i.Logger())
	m := mailer.NewMailer(i.Mailer())

	// Subscribers
	acs := subscriber.NewSubscriber(&cfg.Broker, i.SubAuthCreated(), l)

	// Channels
	ec, err := channels.NewEmailChannel(m, cfg.Client.BaseURL, cfg.Mailer.From)
	if err != nil {
		return nil, err
	}

	// Handlers
	ah := handlers.NewAuthHandler(ec)

	return &Container{
		config: cfg,
		logger: l,
		mailer: m,
		acs:    acs,
		ec:     ec,
		ah:     ah,
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
