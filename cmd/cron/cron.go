package main

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/core"
	"github.com/briskt/go-htmx-app/log"
)

func main() {
	log.Init()

	db, err := app.OpenDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}

	emailSvc, err := app.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to create database transaction: %v", err)
	}

	core.SendPeriodicMessages(context.Background(), tx, emailSvc)

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit database transaction: %v", err)
	}
}
