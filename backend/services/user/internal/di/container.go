package di

import (
	"github.com/ritchieridanko/pasarly/backend/services/user/configs"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/subscriber"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/processors"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/repositories"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger
	acs        *subscriber.Subscriber
	ur         repositories.UserRepository
	up         processors.UserProcessor
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(i.Database())
	tx := database.NewTransactor(i.Database())
	l := logger.NewLogger(i.Logger())

	// Subscribers
	acs := subscriber.NewSubscriber(&cfg.Broker, i.SubAuthCreated(), l)

	// Repositories
	ur := repositories.NewUserRepository(db)

	// Processors
	up := processors.NewUserProcessor(ur, tx)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		acs:        acs,
		ur:         ur,
		up:         up,
	}
}

func (c *Container) SubAuthCreated() *subscriber.Subscriber {
	return c.acs
}

func (c *Container) UserProcessor() processors.UserProcessor {
	return c.up
}
