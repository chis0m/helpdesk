GOOSE_DRIVER ?= mysql
GOOSE_DIR ?= migrations
DB_HOST ?= localhost
DB_DATABASE ?= secure_helpdesk
DB_USERNAME ?= admin
DB_PASSWORD ?= password
DB_PORT ?= 3306
GOOSE_DBSTRING ?= $(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)?parseTime=true

gm: # goose_migration
	goose -dir $(GOOSE_DIR) create $(name) sql

gs: # goose_status
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" status

gu: # goose_up
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up

guo: # goose_up_one
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up-by-one

gd: # goose_down
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down

gr: # goose_down_all
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" reset

dbr: # db_reset — roll back all migrations like `make gr` (local dev only; run from backend/)
	go run ./scripts/db-reset

dbrf: # db_refresh — roll back all then migrate up again,
	go run ./scripts/db-refresh
	
dbm: # db_migrate — goose up (local dev only; run from backend/)
	go run ./scripts/db-migrate

serve: # serve_application
	go run ./cmd/server

check:
	go test ./...
	go fmt ./...

gocheck:
	go vet ./...
	go build -o /dev/null ./cmd/server

govuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

tidy:
	go mod tidy

.PHONY: gm gs gu guo gd gr dbr dbrf dbm serve check gocheck govuln tidy