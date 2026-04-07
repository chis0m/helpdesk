package ca

import (
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
)

const staffDomain = "secweb.ie"

// EmailStaffAdmin and EmailStaffSupport use firstname.lastname@secweb.ie.
var (
	EmailStaffAdmin   = emailAt("Casey", "Admin", staffDomain)
	EmailStaffSupport = emailAt("Sam", "Support", staffDomain)
)

// EnsureStaffUsers seeds one admin and one staff member; both use caTestPassword (same as CA customers).
func EnsureStaffUsers(db *gorm.DB) (admin *models.User, support *models.User, err error) {
	now := time.Now().UTC()
	ptrNow := &now

	hash, err := auth.HashPassword(caTestPassword)
	if err != nil {
		return nil, nil, err
	}

	log := logger.L()

	uAdmin, createdAdmin, err := firstOrCreateUser(
		db,
		EmailStaffAdmin,
		"Casey",
		"Admin",
		models.RoleAdmin,
		false,
		ptrNow,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdAdmin {
		log.Info().
			Str("email", uAdmin.Email).
			Bool("must_change_password", uAdmin.MustChangePassword).
			Msg("CA seed: test user created")
	}

	uSupport, createdSupport, err := firstOrCreateUser(
		db,
		EmailStaffSupport,
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
	if createdSupport {
		log.Info().
			Str("email", uSupport.Email).
			Bool("must_change_password", uSupport.MustChangePassword).
			Msg("CA seed: test user created")
	}

	return uAdmin, uSupport, nil
}
