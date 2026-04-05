package requests

import (
	"time"

	"helpdesk/backend/internal/models"
)

type CreateUserInput struct {
	Email              string
	PasswordHash       string
	FirstName          string
	LastName           string
	MiddleName         *string
	Role               models.UserRole
	IsActive           *bool
	MustChangePassword *bool
	PasswordChangedAt  *time.Time
}

type UpdateUserInput struct {
	Email      *string
	FirstName  *string
	LastName   *string
	MiddleName *string
	IsActive   *bool
}

type UpdateUserRoleRequest struct {
	Role models.UserRole `json:"role" binding:"required,oneof=user staff admin super_admin"`
}

type CreateUserRequest struct {
	Email      string  `json:"email" binding:"required,email,max=120"`
	Password   string  `json:"password" binding:"required,min=8,max=128"`
	FirstName  string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName   string  `json:"last_name" binding:"required,min=2,max=100"`
	MiddleName *string `json:"middle_name" binding:"omitempty,max=100"`
}

type UpdateUserRequest struct {
	Email      *string `json:"email" binding:"omitempty,email,max=120"`
	FirstName  *string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName   *string `json:"last_name" binding:"omitempty,min=2,max=100"`
	MiddleName *string `json:"middle_name" binding:"omitempty,max=100"`
	IsActive   *bool   `json:"is_active" binding:"omitempty"`
}

type CreateStaffRequest struct {
	Email      string  `json:"email" binding:"required,email,max=120"`
	Password   string  `json:"password" binding:"required,min=8,max=128"`
	FirstName  string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName   string  `json:"last_name" binding:"required,min=2,max=100"`
	MiddleName *string `json:"middle_name" binding:"omitempty,max=100"`
	IsActive   *bool   `json:"is_active" binding:"omitempty"`
}
