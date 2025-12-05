package di

import (
	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/interface/server"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/usecases"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/utils"
)

type Container struct {
	config     *configs.Config
	cache      *cache.Cache
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger
	acp        *publisher.Publisher
	ar         repositories.AuthRepository
	sr         repositories.SessionRepository
	tr         repositories.TokenRepository
	bcrypt     *utils.BCrypt
	jwt        *utils.JWT
	validator  *utils.Validator
	au         usecases.AuthUsecase
	su         usecases.SessionUsecase
	ah         *handlers.AuthHandler
	server     *server.Server
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	c := cache.NewCache(&cfg.Cache, i.Cache())
	db := database.NewDatabase(i.Database())
	tx := database.NewTransactor(i.Database())
	l := logger.NewLogger(i.Logger())

	// Publishers
	acp := publisher.NewPublisher(i.PubAuthCreated(), l)

	// Repositories
	ar := repositories.NewAuthRepository(db, c)
	sr := repositories.NewSessionRepository(db)
	tr := repositories.NewTokenRepository(&cfg.Auth, c)

	// Utils
	b := utils.NewBCrypt(cfg.Auth.BCrypt.Cost)
	j := utils.NewJWT(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := utils.NewValidator()

	// Usecases
	au := usecases.NewAuthUsecase(ar, tr, tx, acp, b, v, l)
	su := usecases.NewSessionUsecase(cfg.Auth.Token.Duration.Session, sr, tx, j, v)

	// Handlers
	ah := handlers.NewAuthHandler(au, su, l)

	// Server
	s := server.Init(&cfg.Server, ah, l)

	return &Container{
		config:     cfg,
		cache:      c,
		database:   db,
		transactor: tx,
		logger:     l,
		acp:        acp,
		ar:         ar,
		sr:         sr,
		tr:         tr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		au:         au,
		su:         su,
		ah:         ah,
		server:     s,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
