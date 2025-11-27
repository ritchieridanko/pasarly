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

	ah *handlers.AuthHandler

	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(i.Logger())

	// Utils
	c := utils.NewCookie(cfg, true)

	// Handlers
	ah := handlers.NewAuthHandler(cfg, i.AuthService(), c)

	// Router
	r := router.Init(cfg, ah)

	// Server
	s := server.Init(&cfg.Server, r.Router(), l)

	return &Container{
		config: cfg,
		logger: l,
		cookie: c,
		ah:     ah,
		router: r,
		server: s,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
