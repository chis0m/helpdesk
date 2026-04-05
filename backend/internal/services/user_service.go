package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/repositories"
	"helpdesk/backend/internal/requests"
)

var ErrUserRoleChangeForbidden = errors.New("forbidden role change")

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UpdateUser(userUUID uuid.UUID, input requests.UpdateUserInput) (*models.User, error) {
	return s.userRepo.Update(userUUID, input)
}

func (s *UserService) UpdateRoleByID(userID uint64, role models.UserRole) (*models.User, error) {
	return s.userRepo.UpdateRoleByID(userID, role)
}

func (s *UserService) UpdateRoleByIDAsActor(userID uint64, targetRole models.UserRole, actorRole models.UserRole) (*models.User, error) {
	targetUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	switch actorRole {
	case models.RoleSuperAdmin:
		return s.userRepo.UpdateRoleByID(userID, targetRole)
	case models.RoleAdmin:
		if targetRole != models.RoleUser && targetRole != models.RoleStaff {
			return nil, ErrUserRoleChangeForbidden
		}
		if targetUser.Role == models.RoleAdmin || targetUser.Role == models.RoleSuperAdmin {
			return nil, ErrUserRoleChangeForbidden
		}
		return s.userRepo.UpdateRoleByID(userID, targetRole)
	default:
		return nil, ErrUserRoleChangeForbidden
	}
}

// VULN-02: IDOR on user profiles — loads user by id with no actor-vs-target authorization.
func (s *UserService) GetByID(userID uint64) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) ListAll(page, limit int, role *models.UserRole) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.userRepo.List(page, limit, role)
}

func (s *UserService) GetByUUIDString(userUUID string) (*models.User, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(userUUID))
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetByUUID(parsed)
}

func (s *UserService) CreateStaffFromRequest(req requests.CreateStaffRequest) (*models.User, error) {
	passwordHash, err := auth.HashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	mustChangePassword := true
	input := requests.CreateUserInput{
		Email:              strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash:       passwordHash,
		FirstName:          strings.TrimSpace(req.FirstName),
		LastName:           strings.TrimSpace(req.LastName),
		MiddleName:         req.MiddleName,
		Role:               models.RoleStaff,
		IsActive:           &isActive,
		MustChangePassword: &mustChangePassword,
		PasswordChangedAt:  nil,
	}

	return s.userRepo.Create(input)
}

// VULN-02: IDOR on user profiles — updates user by id with no actor-vs-target authorization.
func (s *UserService) UpdateByID(userID uint64, req requests.UpdateUserRequest) (*models.User, error) {
	input := requests.UpdateUserInput{
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: req.MiddleName,
		IsActive:   req.IsActive,
	}
	return s.userRepo.UpdateByID(userID, input)
}
