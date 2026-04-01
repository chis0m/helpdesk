package routes

import (
	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/container"
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
	}
}
