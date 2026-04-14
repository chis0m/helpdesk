// Local development only: loads .env, runs goose up. Refuses to run unless GO_ENV is development or local.
// Run from the backend module root (e.g. go run ./scripts/db-migrate) so migrations/ resolves.
package main

import (
	"log"
	"strings"

	"helpdesk/backend/internal/config"
	"helpdesk/backend/internal/database"
)

func main() {
	cfg := config.Load()
	if e := strings.ToLower(strings.TrimSpace(cfg.GoEnv)); e != "development" && e != "local" {
		log.Fatalf("db-migrate: refusing to run (GO_ENV=%q; need development or local)", cfg.GoEnv)
	}
	if err := database.RunMigrations(cfg.MySQLDSN()); err != nil {
		log.Fatal(err)
	}
}
