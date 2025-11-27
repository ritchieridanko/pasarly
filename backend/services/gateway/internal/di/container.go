package di

import (
	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/router"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/server"
)

type Container struct {
	logger *logger.Logger
	ah     *handlers.AuthHandler
	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, i *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(i.Logger())

	// Handlers
	ah := handlers.NewAuthHandler(i.AuthService())

	// Router
	r := router.Init(cfg, ah)

	// Server
	s := server.Init(&cfg.Server, r.Router(), l)

	return &Container{
		logger: l,
		ah:     ah,
		router: r,
		server: s,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
