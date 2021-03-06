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
	"github.com/cesarFuhr/gocrypto/internal/pkg/database"
	"github.com/google/uuid"
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

	testdb := bootstrapSQLDatabase(cfg)
	defer testdb.Close()

	err = database.MigrateUp(testdb)
	defer database.MigrateDown(testdb)
	if err != nil {
		panic(err)
	}

	setupDB(testdb)

	httpServer = bootstrapHTTPServer(cfg, testdb)

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
