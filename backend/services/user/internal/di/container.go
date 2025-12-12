package di

import (
	"github.com/ritchieridanko/pasarly/backend/services/user/configs"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/subscriber"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/interface/server"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/processors"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/usecases"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/utils"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger
	acs        *subscriber.Subscriber
	ur         repositories.UserRepository
	ar         repositories.AddressRepository
	up         processors.UserProcessor
	validator  *utils.Validator
	uu         usecases.UserUsecase
	au         usecases.AddressUsecase
	uh         *handlers.UserHandler
	ah         *handlers.AddressHandler
	server     *server.Server
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
	ar := repositories.NewAddressRepository(db)

	// Processors
	up := processors.NewUserProcessor(ur, tx)

	// Utils
	v := utils.NewValidator()

	// Usecases
	uu := usecases.NewUserUsecase(ur, v)
	au := usecases.NewAddressUsecase(ar, tx, v)

	// Handlers
	uh := handlers.NewUserHandler(uu, l)
	ah := handlers.NewAddressHandler(au, l)

	// Server
	s := server.Init(&cfg.Server, l, uh, ah)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		acs:        acs,
		ur:         ur,
		ar:         ar,
		up:         up,
		validator:  v,
		uu:         uu,
		au:         au,
		uh:         uh,
		ah:         ah,
		server:     s,
	}
}

func (c *Container) SubAuthCreated() *subscriber.Subscriber {
	return c.acs
}

func (c *Container) UserProcessor() processors.UserProcessor {
	return c.up
}

func (c *Container) Server() *server.Server {
	return c.server
}
