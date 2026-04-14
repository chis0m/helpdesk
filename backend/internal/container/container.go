package container

import (
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/controllers"
	"helpdesk/backend/internal/mail"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/services"
)

type Container struct {
	DB *gorm.DB

	HealthController    *controllers.HealthController
	AuthController      *controllers.AuthController
	UserController      *controllers.UserController
	InviteController    *controllers.InviteController
	TicketController    *controllers.TicketController
	UserService         *services.UserService
	TicketService       *services.TicketService
	TokenMaker          auth.MakerInterface
	PublicAuthCSRFStore *auth.PublicAuthCSRFStore
	SessionRepo         *repositories.AuthSessionRepository
}

func New(db *gorm.DB, cfg config.Config, tokenMaker auth.MakerInterface) *Container {
	userRepo := repositories.NewUserRepository(db)
	inviteRepo := repositories.NewInviteRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	ticketCommentRepo := repositories.NewTicketCommentRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)
	publicAuthCSRFStore := auth.NewPublicAuthCSRFStore(cfg.CSRFTTL())
	userService := services.NewUserService(userRepo)
	inviteNotifier, resetNotifier := mail.NewNotifiers(cfg)
	inviteService := services.NewInviteService(cfg, inviteRepo, userRepo, inviteNotifier)
	ticketService := services.NewTicketService(ticketRepo, ticketCommentRepo, userRepo)
	authService := services.NewAuthService(cfg, tokenMaker, userRepo, sessionRepo, passwordResetRepo, resetNotifier)
	healthController := controllers.NewHealthController()
	authController := controllers.NewAuthController(cfg, authService, publicAuthCSRFStore, auditLogRepo)
	userController := controllers.NewUserController(userService, auditLogRepo)
	inviteController := controllers.NewInviteController(cfg, inviteService, userService, auditLogRepo)
	ticketController := controllers.NewTicketController(ticketService, userRepo, auditLogRepo, ticketRepo)

	return &Container{
		DB:                  db,
		HealthController:    healthController,
		AuthController:      authController,
		UserController:      userController,
		InviteController:    inviteController,
		TicketController:    ticketController,
		UserService:         userService,
		TicketService:       ticketService,
		TokenMaker:          tokenMaker,
		PublicAuthCSRFStore: publicAuthCSRFStore,
		SessionRepo:         sessionRepo,
	}
}
