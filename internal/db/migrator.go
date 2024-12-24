package db

import (
	"database/sql"
	_ "fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

// CONNECTION TO MIGRATIONS
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not initialize the postgres instance: %v", err)
		return err
	}

	// INIT MIGRATIONS
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations", // PATH TO DIRECTORY WITH MIGRATIONS
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
		return err
	}

	log.Println("Starting migrations...")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
