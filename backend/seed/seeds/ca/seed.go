package ca

import (
	"helpdesk/backend/internal/logger"

	"gorm.io/gorm"
)

// SeedAll runs CA assessment fixtures: customers, staff, tickets. Idempotent.
func SeedAll(db *gorm.DB) error {
	log := logger.L()
	log.Info().Msg("checking CA fixture seed")

	uMust, uOK, err := EnsureCustomerUsers(db)
	if err != nil {
		return err
	}

	_, staffSupport, err := EnsureStaffUsers(db)
	if err != nil {
		return err
	}

	if err := EnsureTickets(db, uMust, uOK, staffSupport); err != nil {
		return err
	}

	log.Info().Msg("CA fixture seed completed")
	return nil
}
