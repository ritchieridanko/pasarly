package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	config *configs.Config
	router *gin.Engine
}

func Init(cfg *configs.Config, l *logger.Logger, ah *handlers.AuthHandler) *Router {
	r := gin.New()
	r.Use(otelgin.Middleware(cfg.App.Name))
	r.Use(gin.Recovery())
	r.Use(middlewares.Logger(l))

	r.ContextWithFallback = true

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "OK",
		})
	})

	v1 := r.Group("/api/v1", middlewares.NewRequestID())

	secret := cfg.JWT.Secret

	// Auth
	auth := v1.Group("/auth")
	{
		auth.POST("/sign-up", ah.SignUp)
		auth.POST("/sign-in", ah.SignIn)
		auth.POST("/sign-out", middlewares.Authenticate(secret), ah.SignOut)
	}

	return &Router{config: cfg, router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
