// Local development only: rolls back all migrations then runs up again (Laravel migrate:refresh).
// Run from the backend module root (e.g. go run ./scripts/db-refresh) so migrations/ resolves.
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
		log.Fatalf("db-refresh: refusing to run (GO_ENV=%q; need development or local)", cfg.GoEnv)
	}
	if err := database.RefreshDb(cfg.MySQLDSN()); err != nil {
		log.Fatal(err)
	}
}
