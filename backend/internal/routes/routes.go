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
		api.GET("/public/app-detail", c.AppDetailController.Get)
		api.GET("/auth/public-csrf-token", c.AuthController.PublicAuthCSRFToken)
		api.POST("/auth/login", middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName), c.AuthController.Login)
		api.POST("/auth/signup", middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName), c.AuthController.Signup)
		api.POST("/auth/forgot-password", middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName), c.AuthController.ForgotPassword)
		api.POST("/auth/reset-password", middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName), c.AuthController.ResetPassword)
		api.GET("/invites/verify", c.InviteController.VerifyInvite)
		api.POST("/invites/accept", middleware.PublicAuthCSRFRequired(c.PublicAuthCSRFStore, auth.CSRFHeaderName), c.InviteController.AcceptInvite)
		api.POST("/auth/refresh", middleware.RefreshTokenRequired(c.TokenMaker, auth.RefreshCookieName), middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.AuthController.Refresh)

		// VULN-06: Insufficient security / audit logging — no audit middleware; sensitive routes use ad-hoc logs only.
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(c.TokenMaker, auth.AccessCookieName), middleware.ActiveSessionRequired(c.SessionRepo))
		{
			protected.GET("/auth/csrf-token", c.AuthController.CSRFToken)
			protected.GET("/auth/me", c.AuthController.Me)
			protected.GET("/auth/sessions", c.AuthController.ListSessions)
			protected.POST("/auth/sessions/revoke-my-other-sessions", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.AuthController.RevokeMyOtherSessions)
			protected.DELETE("/auth/sessions/:session_id", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.AuthController.RevokeSession)
			protected.POST("/auth/logout", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.AuthController.Logout)
			protected.POST("/auth/change-password", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.AuthController.ChangePassword)
			protected.GET("/admin/users", c.UserController.ListAdmin)
			protected.POST("/admin/staff", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.UserController.CreateStaff)
			protected.POST("/admin/invites/staff", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.InviteController.CreateStaffInvite)
			protected.GET("/users/me", c.UserController.GetMe)
			protected.PATCH("/users/me", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.UserController.PatchMe)
			protected.POST("/tickets", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.Create)
			protected.GET("/tickets", c.TicketController.List)
			protected.GET("/tickets/search", c.TicketController.Search)
			protected.GET("/tickets/:id", c.TicketController.GetByID)
			protected.PATCH("/tickets/:id", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.UpdateByID)
			protected.PATCH("/tickets/:id/status", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.UpdateStatus)
			protected.PATCH("/tickets/:id/assign", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.Assign)
			protected.DELETE("/tickets/:id", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.DeleteByID)
			protected.POST("/tickets/:id/comments", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.AddComment)
			protected.GET("/tickets/:id/comments", c.TicketController.ListComments)
			protected.PATCH("/tickets/:id/comments/:commentId", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.UpdateComment)
			protected.DELETE("/tickets/:id/comments/:commentId", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.TicketController.DeleteComment)
			// Role management is restricted by controller/service authorization checks.
			protected.PATCH("/admin/users/:user_id/role", middleware.CSRFRequired(c.SessionRepo, auth.CSRFHeaderName), c.UserController.UpdateRoleByUserID)
		}
	}
}
