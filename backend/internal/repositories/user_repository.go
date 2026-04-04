package repositories

import (
	"errors"
	"strings"

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

func (r *UserRepository) Create(input requests.CreateUserInput) (*models.User, error) {
	role := input.Role
	if role == "" {
		role = models.RoleUser
	}

	user := &models.User{
		UUID:         uuid.New(),
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		MiddleName:   input.MiddleName,
		Role:         role,
		IsActive:     true,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
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
	if input.Role != nil {
		updates["role"] = *input.Role
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
