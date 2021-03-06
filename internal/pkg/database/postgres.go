package database

import (
	"database/sql"
	"fmt"
)

// PGConfigs configuration for a postgres database
type PGConfigs struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	Dbname       string
	MaxOpenConns int
}

// NewPGDatabase Created a connection with the database and returns it
func NewPGDatabase(cfg PGConfigs) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open(cfg.Driver, psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	fmt.Println("Connected to SQLKeyRepository")
	return db, nil
}
