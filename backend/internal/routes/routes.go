package routes

import (
	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/container"
	"helpdesk/backend/internal/middleware"
)

func Register(r *gin.Engine, c *container.Container) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "welcome to secure web helpdesk",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/health", c.HealthController.Ping)
		api.GET("/auth/login-csrf-token", c.AuthController.LoginCSRFToken)
		api.POST(
			"/auth/login",
			middleware.LoginCSRFRequired(c.LoginCSRFStore, auth.CSRFHeaderName),
			c.AuthController.Login,
		)
		api.POST(
			"/auth/refresh",
			middleware.RefreshTokenRequired(c.TokenMaker, auth.RefreshCookieName),
			middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName),
			c.AuthController.Refresh,
		)

		protected := api.Group("")
		protected.Use(
			middleware.AuthRequired(c.TokenMaker, auth.AccessCookieName),
			middleware.ActiveSessionRequired(c.SessionRepo),
		)
		{
			protected.GET("/auth/csrf-token", c.AuthController.CSRFToken)
		}
	}
}
