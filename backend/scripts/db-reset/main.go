// Local development only: loads .env, runs goose down-to-0 then up. Refuses to run unless GO_ENV is development or local.
// Run from the backend module root (e.g. go run ./scripts/db-reset) so migrations/ resolves.
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
		log.Fatalf("db-reset: refusing to run (GO_ENV=%q; need development or local)", cfg.GoEnv)
	}
	if err := database.ResetDb(cfg.MySQLDSN()); err != nil {
		log.Fatal(err)
	}
}
