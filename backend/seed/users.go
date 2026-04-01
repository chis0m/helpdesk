package seed

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/models"
)

func seedAdminUser(db *gorm.DB, cfg config.Config) error {
	email := strings.TrimSpace(cfg.SeedAdminEmail)
	if email == "" {
		return errors.New("SEED_ADMIN_EMAIL cannot be empty")
	}
	if strings.TrimSpace(cfg.SeedAdminPassword) == "" {
		return errors.New("SEED_ADMIN_PASSWORD cannot be empty")
	}

	hash, err := auth.HashPassword(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}

	var middleName *string
	if trimmed := strings.TrimSpace(cfg.SeedAdminMiddleName); trimmed != "" {
		middleName = &trimmed
	}

	attrs := models.User{
		UUID:               uuid.New(),
		PasswordHash:       hash,
		FirstName:          strings.TrimSpace(cfg.SeedAdminFirstName),
		MiddleName:         middleName,
		LastName:           strings.TrimSpace(cfg.SeedAdminLastName),
		Role:               models.RoleAdmin,
		IsActive:           true,
		MustChangePassword: true,
		PasswordChangedAt:  nil,
	}

	admin := models.User{Email: email}
	result := db.Where("email = ?", email).Attrs(attrs).FirstOrCreate(&admin)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		fmt.Printf("[seed] admin user created: %s\n", email)
	} else {
		fmt.Printf("[seed] admin user already exists: %s\n", email)
	}
	return nil
}
