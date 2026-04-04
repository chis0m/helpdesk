package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/container"
	"helpdesk/backend/internal/middleware"
)

func Register(r *gin.Engine, c *container.Container) {
	loginRateLimiter := middleware.NewIPRateLimiter(10, time.Minute)
	signupRateLimiter := middleware.NewIPRateLimiter(5, time.Minute)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "welcome to secure web helpdesk",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/health", c.HealthController.Ping)
		api.GET("/auth/public-csrf-token", c.AuthController.PublicAuthCSRFToken)
		api.POST(
			"/auth/login",
			loginRateLimiter.Middleware(),
			middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName),
			c.AuthController.Login,
		)
		api.POST(
			"/auth/signup",
			signupRateLimiter.Middleware(),
			middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName),
			c.AuthController.Signup,
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
			protected.PATCH(
				"/admin/users/:user_id/role",
				middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName),
				c.UserController.UpdateRoleByUserID,
			)
		}
	}
}
