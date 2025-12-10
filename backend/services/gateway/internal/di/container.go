package di

import (
	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/router"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/server"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/utils"
)

type Container struct {
	config *configs.Config
	logger *logger.Logger
	cookie *utils.Cookie
	ah     *handlers.AuthHandler
	uh     *handlers.UserHandler
	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(i.Logger())

	// Utils
	c := utils.NewCookie(cfg.App.Env, cfg.Server.Host, true)

	// Handlers
	ah := handlers.NewAuthHandler(i.AuthService(), c, cfg.Duration.Session)
	uh := handlers.NewUserHandler(i.UserService())

	// Router
	r := router.Init(l, cfg.App.Name, cfg.JWT.Secret, ah, uh)

	// Server
	s := server.Init(&cfg.Server, r.Router(), l)

	return &Container{
		config: cfg,
		logger: l,
		cookie: c,
		ah:     ah,
		uh:     uh,
		router: r,
		server: s,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
