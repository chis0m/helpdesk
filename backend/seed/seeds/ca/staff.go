package ca

import (
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
)

const staffDomain = "secweb.ie"

var (
	EmailStaffSam    = emailAt("Sam", "Support", staffDomain)
	EmailCasseyAdmin = emailAt("Cassey", "Admin", staffDomain)
)

func EnsureStaffUsers(db *gorm.DB) (sam *models.User, cassey *models.User, err error) {
	now := time.Now().UTC()
	ptrNow := &now

	hash, err := auth.HashPassword(caTestPassword)
	if err != nil {
		return nil, nil, err
	}

	log := logger.L()

	uSam, createdSam, err := firstOrCreateUser(
		db,
		EmailStaffSam,
		"Sam",
		"Support",
		models.RoleStaff,
		false,
		ptrNow,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdSam {
		log.Info().
			Str("email", uSam.Email).
			Bool("must_change_password", uSam.MustChangePassword).
			Msg("CA seed: test user created")
	}

	uCassey, createdCassey, err := firstOrCreateUser(
		db,
		EmailCasseyAdmin,
		"Cassey",
		"Admin",
		models.RoleAdmin,
		false,
		ptrNow,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdCassey {
		log.Info().
			Str("email", uCassey.Email).
			Bool("must_change_password", uCassey.MustChangePassword).
			Msg("CA seed: test user created")
	}
	if err := db.Model(&models.User{}).Where("id = ?", uCassey.ID).Update("role", models.RoleAdmin).Error; err != nil {
		return nil, nil, err
	}

	return uSam, uCassey, nil
}
