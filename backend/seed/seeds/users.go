package seeds

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
)

func SeedAdminUser(db *gorm.DB, cfg config.Config) error {
	log := logger.L()
	log.Info().Msg("checking admin user seed")

	email := strings.TrimSpace(cfg.SeedAdminEmail)
	if email == "" {
		log.Error().Msg("seed admin email is empty")
		return errors.New("SEED_ADMIN_EMAIL cannot be empty")
	}
	if strings.TrimSpace(cfg.SeedAdminPassword) == "" {
		log.Error().Msg("seed admin password is empty")
		return errors.New("SEED_ADMIN_PASSWORD cannot be empty")
	}

	hash, err := auth.HashPassword(cfg.SeedAdminPassword)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash seed admin password")
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
		Role:               models.RoleSuperAdmin,
		IsActive:           true,
		MustChangePassword: true,
		PasswordChangedAt:  nil,
	}

	admin := models.User{Email: email}
	result := db.Where("email = ?", email).Attrs(attrs).FirstOrCreate(&admin)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed creating/finding seed admin user")
		return result.Error
	}
	if admin.Role != models.RoleSuperAdmin {
		if err := db.Model(&admin).Update("role", models.RoleSuperAdmin).Error; err != nil {
			log.Error().Err(err).Str("email", email).Msg("failed to enforce super admin role for seeded user")
			return err
		}
	}
	if result.RowsAffected > 0 {
		log.Info().Str("email", email).Msg("super admin user created")
	} else {
		log.Info().Str("email", email).Msg("super admin user already exists")
	}
	return nil
}
