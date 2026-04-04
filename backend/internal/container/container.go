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

	HealthController *controllers.HealthController
	AuthController   *controllers.AuthController
	UserService      *services.UserService
	TokenMaker       auth.MakerInterface
	LoginCSRFStore   *auth.LoginCSRFStore
	SessionRepo      *repositories.AuthSessionRepository
}

func New(db *gorm.DB, cfg config.Config, tokenMaker auth.MakerInterface) *Container {
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)
	loginCSRFStore := auth.NewLoginCSRFStore(cfg.CSRFTTL())
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(cfg, tokenMaker, userRepo, sessionRepo)
	healthController := controllers.NewHealthController()
	authController := controllers.NewAuthController(cfg, authService, loginCSRFStore)

	return &Container{
		DB:               db,
		HealthController: healthController,
		AuthController:   authController,
		UserService:      userService,
		TokenMaker:       tokenMaker,
		LoginCSRFStore:   loginCSRFStore,
		SessionRepo:      sessionRepo,
	}
}
