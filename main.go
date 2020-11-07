package main

import (
	"log"
	"net/http"

	"github.com/cesarFuhr/gocrypto/keys"
)

var syncKeySource = keys.SynchronousKeySource{}
var inMemKeyRepo = keys.InMemoryKeyRepository{Store: make(map[string]keys.Key)}
var keyStore = keys.KeyStore{Source: &syncKeySource, Repo: &inMemKeyRepo}
var crypto = JWECrypto{}

func main() {
	server := &KeyServer{&keyStore, &crypto}
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
