package database

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open migration db: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func ResetDb(dsn string) error {
	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open migration db: %w", err)
	}
	defer db.Close()

	if err := goose.DownTo(db, "migrations", 0); err != nil {
		return fmt.Errorf("goose down-to 0: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func init() {
	goose.SetBaseFS(os.DirFS("."))
}
