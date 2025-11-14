package di

import (
	"github.com/ritchieridanko/pasarly/auth-service/configs"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/repositories"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/usecases"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/cache"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/database"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/publisher"
	"github.com/ritchieridanko/pasarly/auth-service/internal/interface/grpc/handlers"
	"github.com/ritchieridanko/pasarly/auth-service/internal/interface/grpc/server"
	"github.com/ritchieridanko/pasarly/auth-service/internal/service"
)

type Container struct {
	logger     *logger.Logger
	cache      *cache.Cache
	database   *database.Database
	transactor *database.Transactor

	ar repositories.AuthRepository
	sr repositories.SessionRepository
	tr repositories.TokenRepository

	au usecases.AuthUsecase
	su usecases.SessionUsecase

	agh *handlers.AuthGRPCHandler

	gs *server.GRPCServer

	service *service.Service
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(i.Logger())
	c := cache.NewCache(&cfg.Cache, i.Cache())
	db := database.NewDatabase(i.DB())
	tx := database.NewTransactor(i.DB())
	acp := publisher.NewPublisher(i.PubAuthCreated(), l)

	// Repositories
	ar := repositories.NewAuthRepository(db, c)
	sr := repositories.NewSessionRepository(db)
	tr := repositories.NewTokenRepository(&cfg.Auth, c)

	// Service
	s := service.Init(cfg)

	// Usecases
	au := usecases.NewAuthUsecase(ar, tr, tx, acp, s.BCrypt(), s.Validator(), l)
	su := usecases.NewSessionUsecase(&cfg.Auth, sr, tx, s.JWT())

	// Handlers
	agh := handlers.NewAuthGRPCHandler(au, su, l)

	// Servers
	gs := server.NewGRPCServer(&cfg.Server, agh, l)

	return &Container{
		logger:     l,
		cache:      c,
		database:   db,
		transactor: tx,
		ar:         ar,
		sr:         sr,
		tr:         tr,
		au:         au,
		su:         su,
		agh:        agh,
		gs:         gs,
		service:    s,
	}
}

func (c *Container) GRPCServer() *server.GRPCServer {
	return c.gs
}
