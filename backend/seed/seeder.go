package seed

import (
	"fmt"

	"helpdesk/backend/internal/config"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB, cfg config.Config) error {
	fmt.Println("[seed] starting seed process")
	if err := seedAdminUser(db, cfg); err != nil {
		return err
	}
	fmt.Println("[seed] seeding completed successfully")
	return nil
}
