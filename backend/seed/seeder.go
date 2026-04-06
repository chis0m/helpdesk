package seed

import (
	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/seed/seeds"
	"helpdesk/backend/seed/seeds/ca"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB, cfg config.Config) error {
	log := logger.L()
	log.Info().Msg("starting seed process")
	if err := seeds.SeedAdminUser(db, cfg); err != nil {
		log.Error().Err(err).Msg("seeding failed")
		return err
	}
	if cfg.SeedCA {
		if err := ca.SeedAll(db); err != nil {
			log.Error().Err(err).Msg("CA fixture seeding failed")
			return err
		}
	} else {
		log.Info().Msg("CA fixture seed skipped (SEED_CA=false)")
	}
	log.Info().Msg("seeding completed successfully")
	return nil
}
