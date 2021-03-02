package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/cesarFuhr/gocrypto/internal/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"

	// loads the file driver to migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var httpServer *http.Server

func TestMain(m *testing.M) {
	os.Exit(deferable(m))
}

func deferable(m *testing.M) int {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	db := bootstrapSQLDatabase(cfg)
	defer db.Close()

	err = runMigrationsUp(db)
	defer runMigrationsDown(db)
	if err != nil {
		panic(err)
	}

	setupDB(db)

	httpServer = bootstrapHTTPServer(cfg, db)

	return m.Run()
}

func setupDB(db *sql.DB) {
	stmt := `INSERT INTO keys (id, scope, expiration, creation, priv, pub)
	VALUES ($1, $2, $3, $4, $5, $6)`
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	k := keys.Key{
		ID:         uuid.New().String(),
		Scope:      "test",
		Expiration: time.Now().AddDate(0, 0, 1),
		Priv:       rsaKey,
		Pub:        &rsaKey.PublicKey,
	}
	mockKey = k

	_, err := db.Exec(
		stmt,
		k.ID,
		k.Scope,
		k.Expiration,
		time.Now(),
		x509.MarshalPKCS1PrivateKey(k.Priv),
		x509.MarshalPKCS1PublicKey(k.Pub),
	)
	if err != nil {
		panic(err)
	}
}

func runMigrationsUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../internal/pkg/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			return nil
		}
		panic(err)
	}
	return nil
}

func runMigrationsDown(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://../internal/pkg/db/migrations",
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
