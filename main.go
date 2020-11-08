package main

import (
	"crypto/rsa"
	"log"
	"net/http"

	"github.com/cesarFuhr/gocrypto/keys"
)

var pool = make(chan *rsa.PrivateKey, 10)
var gen = keys.KeyGenerator{}
var poolKeySource = keys.PoolKeySource{Pool: pool, Kgen: &gen}
var inMemKeyRepo = keys.InMemoryKeyRepository{Store: make(map[string]keys.Key)}
var keyStore = keys.KeyStore{Source: &poolKeySource, Repo: &inMemKeyRepo}
var crypto = JWECrypto{}

func main() {
	run()
}

func run() {
	poolKeySource.WarmUp()
	server := &KeyServer{&keyStore, &crypto}
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
