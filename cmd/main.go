package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	server "github.com/cesarFuhr/gocrypto/internal/app"
	"github.com/cesarFuhr/gocrypto/internal/app/adapters"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/crypto"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/cesarFuhr/gocrypto/internal/app/ports"
	"github.com/cesarFuhr/gocrypto/internal/pkg/config"
	"github.com/cesarFuhr/gocrypto/internal/pkg/db"
	"github.com/cesarFuhr/gocrypto/internal/pkg/logger"
)

func main() {
	run()
}

func run() {

	cfgSource := getCfgSource()
	cfg, err := config.LoadConfigs(cfgSource)
	if err != nil {
		panic(err)
	}

	db := bootstrapSQLDatabase(cfg)
	httpServer := bootstrapHTTPServer(cfg, db)

	if err := http.ListenAndServe(":"+cfg.Server.Port, httpServer); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

func bootstrapSQLDatabase(cfg config.Config) *sql.DB {
	sqlDB, err := db.NewPGDatabase(db.PGConfigs{
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		User:     cfg.Db.User,
		Password: cfg.Db.Password,
		Dbname:   cfg.Db.Dbname,
		Driver:   cfg.Db.Driver,
	})
	if err != nil {
		panic(err)
	}
	return sqlDB
}

func bootstrapHTTPServer(cfg config.Config, sqlDB *sql.DB) server.HTTPServer {
	keySource := adapters.NewPoolKeySource(cfg.App.KeySource.RSAKeySize, cfg.App.KeySource.PoolSize)
	keySource.WarmUp()

	sqlKeyRepo := adapters.NewSQLKeyRepository(sqlDB)

	keyService := keys.NewKeyService(&keySource, &sqlKeyRepo)
	keyHandler := ports.NewKeyHandler(keyService)

	cryptoService := crypto.NewCryptoService(&sqlKeyRepo)
	encryptHandler := ports.NewEncryptHandler(cryptoService)
	decryptHandler := ports.NewDecryptHandler(cryptoService)

	logger := logger.NewLogger()

	return server.NewHTTPServer(logger, keyHandler, encryptHandler, decryptHandler)
}

func getCfgSource() string {
	var cfgFromEnv bool
	flag.BoolVar(&cfgFromEnv, "e", false, "load config from environment")
	flag.Parse()
	if cfgFromEnv == true {
		return "env"
	}
	return "yaml"
}
