package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// loads the file driver to migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp runs the migrations to setup the db
func MigrateUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Printf("Migrations: no change")
			return nil
		}
		panic(err)
	}
	return nil
}

// MigrateDown runs the migrations to teardown the db
func MigrateDown(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Down(); err != nil {
		if err.Error() == "no change" {
			return nil
		}
		panic(err)
	}
	return nil
}
