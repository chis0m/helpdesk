package container

import (
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/controllers"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/services"
)

type Container struct {
	DB *gorm.DB

	HealthController *controllers.HealthController
	UserService      *services.UserService
	TokenMaker       auth.MakerInterface
	SessionRepo      *repositories.AuthSessionRepository
}

func New(db *gorm.DB, tokenMaker auth.MakerInterface) *Container {
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)
	userService := services.NewUserService(userRepo)
	healthController := controllers.NewHealthController()

	return &Container{
		DB:               db,
		HealthController: healthController,
		UserService:      userService,
		TokenMaker:       tokenMaker,
		SessionRepo:      sessionRepo,
	}
}
