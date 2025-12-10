package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/handlers"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func Init(l *logger.Logger, appName, jwtSecret string, ah *handlers.AuthHandler, uh *handlers.UserHandler) *Router {
	r := gin.New()
	r.Use(otelgin.Middleware(appName))
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

	// Auth
	auth := v1.Group("/auth")
	{
		auth.GET("/email/available", ah.IsEmailAvailable)
		auth.POST("/sign-up", ah.SignUp)
		auth.POST("/sign-in", ah.SignIn)
		auth.POST("/sign-out", middlewares.Authenticate(jwtSecret), ah.SignOut)
	}

	// Users
	users := v1.Group("/users")
	{
		users.GET(
			"/me",
			middlewares.Authenticate(jwtSecret),
			middlewares.Authorize(constants.RoleCustomer),
			uh.GetUser,
		)

		users.PUT(
			"/me",
			middlewares.Authenticate(jwtSecret),
			middlewares.Authorize(constants.RoleCustomer),
			uh.UpsertUser,
		)

		users.PATCH(
			"/me",
			middlewares.Authenticate(jwtSecret),
			middlewares.Authorize(constants.RoleCustomer),
			uh.UpdateUser,
		)
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
