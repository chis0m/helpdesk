package ca

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
)

const caTestPassword = "password"

func emailAt(firstName, lastName, domain string) string {
	local := strings.ToLower(strings.TrimSpace(firstName)) + "." + strings.ToLower(strings.TrimSpace(lastName))
	return local + "@" + strings.TrimSpace(domain)
}

func firstOrCreateUser(
	db *gorm.DB,
	email, firstName, lastName string,
	role models.UserRole,
	mustChange bool,
	changedAt *time.Time,
	passwordHash string,
) (*models.User, bool, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	var found []models.User
	if err := db.Where("email = ?", normalized).Limit(1).Find(&found).Error; err != nil {
		return nil, false, err
	}
	if len(found) > 0 {
		return &found[0], false, nil
	}

	u := models.User{
		UUID:               uuid.New(),
		Email:              normalized,
		PasswordHash:       passwordHash,
		FirstName:          strings.TrimSpace(firstName),
		LastName:           strings.TrimSpace(lastName),
		Role:               role,
		IsActive:           true,
		MustChangePassword: mustChange,
		PasswordChangedAt:  changedAt,
	}
	if err := db.Create(&u).Error; err != nil {
		return nil, false, err
	}
	return &u, true, nil
}
