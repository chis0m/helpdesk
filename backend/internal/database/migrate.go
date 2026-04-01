package database

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	fmt.Println("[migration] starting goose up")

	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open migration db: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	fmt.Println("[migration] goose up completed successfully")

	return nil
}

func ResetDb(dsn string) error {
	fmt.Println("[migration] resetting database (down-to 0 then up)")

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
	fmt.Println("[migration] database reset completed successfully")

	return nil
}

func init() {
	goose.SetBaseFS(os.DirFS("."))
}
