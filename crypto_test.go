package main

import (
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/keys"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

var key = keys.Key{
	Scope:      "scope",
	ID:         "id",
	Expiration: time.Now().AddDate(0, 0, 1),
	Priv:       rsaKey,
	Pub:        &rsaKey.PublicKey,
}

func TestCryptoEncrypt(t *testing.T) {
	crypto := JWECrypto{}
	t.Run("Should return a valid JWE", func(t *testing.T) {
		got, _ := crypto.Encrypt(key.Pub, "testingOK")

		if _, err := jwe.Decrypt(got, jwa.RSA_OAEP_256, key.Priv); err != nil {
			t.Errorf("Invalid jwe: %v", err)
		}
	})
}
