package seed

import (
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/seed/seeds"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB, cfg config.Config) error {
	log := logger.L()
	log.Info().Msg("starting seed process")
	if err := seeds.SeedAdminUser(db, cfg); err != nil {
		log.Error().Err(err).Msg("seeding failed")
		return err
	}
	log.Info().Msg("seeding completed successfully")
	return nil
}
