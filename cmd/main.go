package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	server "github.com/cesarFuhr/gocrypto/internal/app"
	"github.com/cesarFuhr/gocrypto/internal/app/adapters"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/crypto"
	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	"github.com/cesarFuhr/gocrypto/internal/app/ports"
	"github.com/cesarFuhr/gocrypto/internal/pkg/config"
	"github.com/cesarFuhr/gocrypto/internal/pkg/database"
	"github.com/cesarFuhr/gocrypto/internal/pkg/exit"
	"github.com/cesarFuhr/gocrypto/internal/pkg/logger"
)

func main() {
	run()
}

func run() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	db := bootstrapSQLDatabase(cfg)
	database.MigrateUp(db)

	httpServer := bootstrapHTTPServer(cfg, db)

	e := make(chan struct{}, 1)
	exit.ListenToExit(e)

	go gracefullShutdown(e, httpServer)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

func bootstrapSQLDatabase(cfg config.Config) *sql.DB {
	sqlDB, err := database.NewPGDatabase(database.PGConfigs{
		Host:         cfg.Db.Host,
		Port:         cfg.Db.Port,
		User:         cfg.Db.User,
		Password:     cfg.Db.Password,
		Dbname:       cfg.Db.Dbname,
		Driver:       cfg.Db.Driver,
		MaxOpenConns: cfg.Db.MaxOpenConns,
	})
	if err != nil {
		panic(err)
	}
	return sqlDB
}

func bootstrapHTTPServer(cfg config.Config, sqlDB *sql.DB) *http.Server {
	keySource := adapters.NewPoolKeySource(cfg.App.KeySource.RSAKeySize, cfg.App.KeySource.PoolSize)
	keySource.WarmUp()

	sqlKeyRepo := adapters.NewSQLKeyRepository(sqlDB)

	keyService := keys.NewKeyService(&keySource, &sqlKeyRepo)
	keyHandler := ports.NewKeyHandler(keyService)

	cryptoService := crypto.NewCryptoService(&sqlKeyRepo)
	encryptHandler := ports.NewEncryptHandler(cryptoService)
	decryptHandler := ports.NewDecryptHandler(cryptoService)

	logger := logger.NewLogger()

	s := server.NewHTTPServer(logger, keyHandler, encryptHandler, decryptHandler)
	s.Addr = ":" + cfg.Server.Port

	return s
}

func gracefullShutdown(e chan struct{}, s *http.Server) {
	<-e
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown properly...")
	}

	cancel()
}
