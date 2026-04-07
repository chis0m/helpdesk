package ca

import (
	"time"

	"gorm.io/gorm"

	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/internal/models"
)

const staffDomain = "secweb.ie"

// EmailStaffSam and EmailStaffCassey use firstname.lastname@secweb.ie (both staff).
var (
	EmailStaffSam    = emailAt("Sam", "Support", staffDomain)
	EmailStaffCassey = emailAt("Cassey", "Support", staffDomain)
)

// EnsureStaffUsers seeds two staff members (Sam + Cassey); both use caTestPassword (same as CA customers).
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
		EmailStaffCassey,
		"Cassey",
		"Support",
		models.RoleStaff,
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

	return uSam, uCassey, nil
}
