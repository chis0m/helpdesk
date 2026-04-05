package container

import (
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/controllers"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/services"
)

type Container struct {
	DB *gorm.DB

	HealthController    *controllers.HealthController
	AuthController      *controllers.AuthController
	UserController      *controllers.UserController
	TicketController    *controllers.TicketController
	UserService         *services.UserService
	TicketService       *services.TicketService
	TokenMaker          auth.MakerInterface
	PublicAuthCSRFStore *auth.PublicAuthCSRFStore
	SessionRepo         *repositories.AuthSessionRepository
}

func New(db *gorm.DB, cfg config.Config, tokenMaker auth.MakerInterface) *Container {
	userRepo := repositories.NewUserRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)
	publicAuthCSRFStore := auth.NewPublicAuthCSRFStore(cfg.CSRFTTL())
	userService := services.NewUserService(userRepo)
	ticketService := services.NewTicketService(ticketRepo, userRepo)
	authService := services.NewAuthService(cfg, tokenMaker, userRepo, sessionRepo)
	healthController := controllers.NewHealthController()
	authController := controllers.NewAuthController(cfg, authService, publicAuthCSRFStore)
	userController := controllers.NewUserController(userService)
	ticketController := controllers.NewTicketController(ticketService)

	return &Container{
		DB:                  db,
		HealthController:    healthController,
		AuthController:      authController,
		UserController:      userController,
		TicketController:    ticketController,
		UserService:         userService,
		TicketService:       ticketService,
		TokenMaker:          tokenMaker,
		PublicAuthCSRFStore: publicAuthCSRFStore,
		SessionRepo:         sessionRepo,
	}
}
