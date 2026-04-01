package main

import (
	"log"

	"helpdesk/backend/boot"
)

func main() {
	app, err := boot.NewApp()
	if err != nil {
		log.Fatalf("failed to bootstrap app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
