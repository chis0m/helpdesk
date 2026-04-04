package requests

import (
	"time"

	"helpdesk/backend/internal/models"
)

type CreateUserInput struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	MiddleName   *string
	Role         models.UserRole
	IsActive           *bool
	MustChangePassword *bool
	PasswordChangedAt  *time.Time
}

type UpdateUserInput struct {
	Email      *string
	FirstName  *string
	LastName   *string
	MiddleName *string
	Role       *models.UserRole
	IsActive   *bool
}

type UpdateUserRoleRequest struct {
	Role models.UserRole `json:"role" binding:"required,oneof=user staff admin"`
}
