// Local development only: rolls back all migrations then runs up again (Laravel migrate:refresh).
// When SEED_CA is true, also runs seed.SeedAll.
// Run from the backend module root (e.g. go run ./scripts/db-refresh) so migrations/ resolves.
package main

import (
	"log"
	"strings"

	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/database"
	"helpdesk/backend/internal/logger"
	"helpdesk/backend/seed"
)

func main() {
	cfg := config.Load()
	if e := strings.ToLower(strings.TrimSpace(cfg.GoEnv)); e != "development" && e != "local" {
		log.Fatalf("db-refresh: refusing to run (GO_ENV=%q; need development or local)", cfg.GoEnv)
	}
	if err := database.RefreshDb(cfg.MySQLDSN()); err != nil {
		log.Fatal(err)
	}
	if !cfg.SeedCA {
		log.Println("db-refresh: seed skipped (SEED_CA=false)")
		return
	}
	if err := logger.Init(cfg); err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sdb, _ := db.DB()
		if sdb != nil {
			_ = sdb.Close()
		}
	}()
	if err := seed.SeedAll(db, cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("db-refresh: seed completed")
}
