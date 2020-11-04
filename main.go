package main

import (
	"log"
	"net/http"

	"github.com/cesarFuhr/gocrypto/keys"
)

var inMemKeySource = keys.InMemoryKeySource{}
var inMemKeyRepo = keys.InMemoryKeyRepository{Store: make(map[string]keys.Key)}
var keyStore = keys.KeyStore{Source: &inMemKeySource, Repo: &inMemKeyRepo}

func main() {
	server := &KeyServer{&keyStore, nil}
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
