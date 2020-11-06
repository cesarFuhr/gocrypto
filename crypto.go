package main

import (
	"crypto/rsa"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

type JWECrypto struct{}

// Encrypt Encrypts the content in a JWE Wrapper
func (c JWECrypto) Encrypt(k *rsa.PublicKey, m string) ([]byte, error) {
	msg, _ := jwe.Encrypt([]byte(m), jwa.RSA_OAEP_256, k, jwa.A256CBC_HS512, jwa.NoCompress)
	return msg, nil
}
