package container

import (
	"gorm.io/gorm"

	"helpdesk/backend/internal/controllers"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/services"
)

type Container struct {
	DB *gorm.DB

	HealthController *controllers.HealthController
	UserService      *services.UserService
}

func New(db *gorm.DB) *Container {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	healthController := controllers.NewHealthController()

	return &Container{
		DB:               db,
		HealthController: healthController,
		UserService:      userService,
	}
}
