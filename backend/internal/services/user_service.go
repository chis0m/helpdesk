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

var (
	ErrUserRoleChangeForbidden         = errors.New("forbidden role change")
	ErrCreateStaffAdminRequiresSuperAdmin = errors.New("only super_admin may create staff with role admin")
)

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

	// Only super_admin may change roles. Staff may only be promoted to admin or super_admin
	// (never to customer/user). Customers (user) must never be given the staff role via this API.
	if actorRole != models.RoleSuperAdmin {
		return nil, ErrUserRoleChangeForbidden
	}
	if err := validateRoleTransition(targetUser.Role, targetRole); err != nil {
		return nil, err
	}
	return s.userRepo.UpdateRoleByID(userID, targetRole)
}

func validateRoleTransition(current, target models.UserRole) error {
	if current == target {
		return nil
	}
	if current == models.RoleUser && target == models.RoleStaff {
		return ErrUserRoleChangeForbidden
	}
	if current == models.RoleStaff && target == models.RoleUser {
		return ErrUserRoleChangeForbidden
	}
	if current == models.RoleStaff && target != models.RoleAdmin && target != models.RoleSuperAdmin {
		return ErrUserRoleChangeForbidden
	}
	return nil
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

func (s *UserService) CreateStaffFromRequest(actorRole models.UserRole, req requests.CreateStaffRequest) (*models.User, error) {
	targetRole := models.RoleStaff
	if r := strings.TrimSpace(req.Role); r != "" {
		targetRole = models.UserRole(r)
	}
	if targetRole == models.RoleAdmin && actorRole != models.RoleSuperAdmin {
		return nil, ErrCreateStaffAdminRequiresSuperAdmin
	}

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
		Role:               targetRole,
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
