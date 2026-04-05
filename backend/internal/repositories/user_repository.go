package repositories

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/requests"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	normalized := strings.ToLower(strings.TrimSpace(email))
	if err := r.db.First(&user, "email = ?", normalized).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUUID(userUUID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(userID uint64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(input requests.CreateUserInput) (*models.User, error) {
	role := input.Role
	if role == "" {
		role = models.RoleUser
	}
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}
	mustChangePassword := true
	if input.MustChangePassword != nil {
		mustChangePassword = *input.MustChangePassword
	}

	user := &models.User{
		UUID:               uuid.New(),
		Email:              input.Email,
		PasswordHash:       input.PasswordHash,
		FirstName:          input.FirstName,
		LastName:           input.LastName,
		MiddleName:         input.MiddleName,
		Role:               role,
		IsActive:           isActive,
		MustChangePassword: mustChangePassword,
		PasswordChangedAt:  input.PasswordChangedAt,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateRoleByID(userID uint64, role models.UserRole) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&user).Update("role", role).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user role updated but failed to reload")
	}

	return &user, nil
}

func (r *UserRepository) UpdateByID(userID uint64, input requests.UpdateUserInput) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.FirstName != nil {
		updates["first_name"] = *input.FirstName
	}
	if input.LastName != nil {
		updates["last_name"] = *input.LastName
	}
	if input.MiddleName != nil {
		updates["middle_name"] = *input.MiddleName
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) == 0 {
		return &user, nil
	}

	if err := r.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user updated but failed to reload")
	}

	return &user, nil
}

func (r *UserRepository) UpdatePasswordByUUID(userUUID uuid.UUID, passwordHash string, mustChangePassword bool, changedAt time.Time) error {
	return r.db.Model(&models.User{}).
		Where("uuid = ?", userUUID).
		Updates(map[string]any{
			"password_hash":        passwordHash,
			"must_change_password": mustChangePassword,
			"password_changed_at":  changedAt.UTC(),
		}).Error
}

func (r *UserRepository) Update(userUUID uuid.UUID, input requests.UpdateUserInput) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.FirstName != nil {
		updates["first_name"] = *input.FirstName
	}
	if input.LastName != nil {
		updates["last_name"] = *input.LastName
	}
	if input.MiddleName != nil {
		updates["middle_name"] = *input.MiddleName
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) == 0 {
		return &user, nil
	}

	if err := r.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&user, "uuid = ?", userUUID).Error; err != nil {
		return nil, errors.New("user updated but failed to reload")
	}

	return &user, nil
}
