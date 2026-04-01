package requests

import "helpdesk/backend/internal/models"

type CreateUserInput struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	MiddleName   *string
	Role         models.UserRole
}

type UpdateUserInput struct {
	Email      *string
	FirstName  *string
	LastName   *string
	MiddleName *string
	Role       *models.UserRole
	IsActive   *bool
}
