GOOSE_DRIVER ?= mysql
GOOSE_DIR ?= migrations
DB_HOST ?= localhost
DB_DATABASE ?= helpdesk
DB_USERNAME ?= admin
DB_PASSWORD ?= password
DB_PORT ?= 3306
GOOSE_DBSTRING ?= $(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)?parseTime=true

goose_operations: gm gs gu guo gd gr

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

serve: # serve_application
	go run ./cmd/server

check:
	go test ./...
	go fmt ./...

tidy:
	go mod tidy

.PHONY: goose_operations serve check tidy