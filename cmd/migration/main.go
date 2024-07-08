package main

import (
	"log"

	"github.com/mbatimel/WB_Weather_Service/internal/migrate"
	"github.com/mbatimel/WB_Weather_Service/internal/repo"
)

func main() {
	db, err := repo.SetConfigs("config/config.yaml")
    if err != nil {
        log.Fatalf("Error setting up database: %v", err)
    }
    defer db.Close()

    err = db.ConnectToDataBase()
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    // Example migration application
    err = migrate.ApplyMigrations(db, "migrations/migrate.sql")
    if err != nil {
        log.Fatalf("Error applying migrations: %v", err)
    }

    log.Println("Migrations applied successfully!")

}