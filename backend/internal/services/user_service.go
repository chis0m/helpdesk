package services

import (
	"github.com/google/uuid"

	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(input requests.CreateUserInput) (*models.User, error) {
	return s.userRepo.Create(input)
}

func (s *UserService) UpdateUser(userUUID uuid.UUID, input requests.UpdateUserInput) (*models.User, error) {
	return s.userRepo.Update(userUUID, input)
}
