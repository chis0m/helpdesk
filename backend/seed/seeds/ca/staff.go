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

// EnsureStaffUsers seeds one admin and one staff member; both receive random passwords logged on first creation only.
func EnsureStaffUsers(db *gorm.DB) (admin *models.User, support *models.User, err error) {
	log := logger.L()
	now := time.Now().UTC()
	ptrNow := &now

	plainAdmin, err := randomPassword()
	if err != nil {
		return nil, nil, err
	}
	hashAdmin, err := auth.HashPassword(plainAdmin)
	if err != nil {
		return nil, nil, err
	}
	uAdmin, createdAdmin, err := firstOrCreateUser(
		db,
		EmailStaffAdmin,
		"Casey",
		"Admin",
		models.RoleAdmin,
		false,
		ptrNow,
		hashAdmin,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdAdmin {
		log.Warn().
			Str("email", uAdmin.Email).
			Str("plaintext_password", plainAdmin).
			Msg("CA seed: random password generated on purpose for test data")
	}

	plainSupport, err := randomPassword()
	if err != nil {
		return nil, nil, err
	}
	hashSupport, err := auth.HashPassword(plainSupport)
	if err != nil {
		return nil, nil, err
	}
	uSupport, createdSupport, err := firstOrCreateUser(
		db,
		EmailStaffSupport,
		"Sam",
		"Support",
		models.RoleStaff,
		false,
		ptrNow,
		hashSupport,
	)
	if err != nil {
		return nil, nil, err
	}
	if createdSupport {
		log.Warn().
			Str("email", uSupport.Email).
			Str("plaintext_password", plainSupport).
			Msg("CA seed: random password generated on purpose for test data")
	}

	return uAdmin, uSupport, nil
}
