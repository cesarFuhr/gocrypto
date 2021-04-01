package database

import (
	"database/sql"
	"embed"
	"io/fs"
)

//go:embed migrations
var mFS embed.FS

// MigrateUp runs the migrations to setup the db
func MigrateUp(db *sql.DB) error {
	upFilePaths, err := fs.Glob(mFS, "migrations/*.up.sql")
	if err != nil {
		return err
	}

	for _, fp := range upFilePaths {
		m, err := fs.ReadFile(mFS, fp)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(m) + ";")
		if err != nil {
			return err
		}
	}

	return nil
}

// MigrateDown runs the migrations to teardown the db
func MigrateDown(db *sql.DB) error {
	downFilePaths, err := fs.Glob(mFS, "migrations/*.down.sql")
	if err != nil {
		return err
	}

	for _, fp := range downFilePaths {
		m, err := fs.ReadFile(mFS, fp)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(m) + ";")
		if err != nil {
			return err
		}
	}

	return nil
}
