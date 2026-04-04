package services

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"helpdesk/backend/internal/auth"
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

func (s *UserService) UpdateRoleByID(userID uint64, role models.UserRole) (*models.User, error) {
	return s.userRepo.UpdateRoleByID(userID, role)
}

func (s *UserService) GetByID(userID uint64) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) CreateUserFromRequest(req requests.CreateUserRequest) (*models.User, error) {
	passwordHash, err := auth.HashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		return nil, err
	}

	mustChangePassword := false
	changedAt := time.Now().UTC()
	input := requests.CreateUserInput{
		Email:              strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash:       passwordHash,
		FirstName:          strings.TrimSpace(req.FirstName),
		LastName:           strings.TrimSpace(req.LastName),
		MiddleName:         req.MiddleName,
		Role:               req.Role,
		IsActive:           req.IsActive,
		MustChangePassword: &mustChangePassword,
		PasswordChangedAt:  &changedAt,
	}

	return s.userRepo.Create(input)
}

func (s *UserService) UpdateByID(userID uint64, req requests.UpdateUserRequest) (*models.User, error) {
	input := requests.UpdateUserInput{
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: req.MiddleName,
		Role:       req.Role,
		IsActive:   req.IsActive,
	}
	return s.userRepo.UpdateByID(userID, input)
}
