package database

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"

	"helpdesk/backend/internal/logger"
)

func RunMigrations(dsn string) error {
	log := logger.L()
	log.Info().Msg("starting goose up migrations")

	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		log.Error().Err(err).Msg("failed opening migration db")
		return fmt.Errorf("open migration db: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		log.Error().Err(err).Msg("goose up failed")
		return fmt.Errorf("goose up: %w", err)
	}
	log.Info().Msg("goose up completed successfully")

	return nil
}

func ResetDb(dsn string) error {
	log := logger.L()
	log.Info().Str("operation", "down-to-0-then-up").Msg("resetting database")

	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		log.Error().Err(err).Msg("failed opening migration db")
		return fmt.Errorf("open migration db: %w", err)
	}
	defer db.Close()

	if err := goose.DownTo(db, "migrations", 0); err != nil {
		log.Error().Err(err).Msg("goose down-to 0 failed")
		return fmt.Errorf("goose down-to 0: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Error().Err(err).Msg("goose up failed")
		return fmt.Errorf("goose up: %w", err)
	}
	log.Info().Msg("database reset completed successfully")

	return nil
}

func init() {
	goose.SetBaseFS(os.DirFS("."))
}
