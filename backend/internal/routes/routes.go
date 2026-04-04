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
		api.POST("/auth/login", c.AuthController.Login)
		api.POST("/auth/refresh", middleware.RefreshTokenRequired(c.TokenMaker, auth.RefreshCookieName), c.AuthController.Refresh)
	}
}
