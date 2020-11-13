package main

import (
	"crypto/rsa"
	"log"
	"net/http"

	"github.com/cesarFuhr/gocrypto/keys"
)

var pool = make(chan *rsa.PrivateKey, 10)
var store = map[string]keys.Key{}
var gen = keys.KeyGenerator{}
var poolKeySource = keys.PoolKeySource{Pool: pool, Kgen: &gen}
var sqlKeyRepo = keys.SQLKeyRepository{Cfg: keys.SQLConfigs{
	Host:     "db",
	Port:     5432,
	User:     "postgres",
	Password: "pass",
	Dbname:   "gocrypto",
	Driver:   "postgres",
}}
var keyStore = keys.KeyStore{Source: &poolKeySource, Repo: &sqlKeyRepo}
var crypto = JWECrypto{}

func main() {
	run()
}

func run() {
	poolKeySource.WarmUp()
	sqlKeyRepo.Connect()
	server := &KeyServer{&keyStore, &crypto}
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
