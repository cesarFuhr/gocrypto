package crypto

import (
	"crypto/rsa"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

// JWECrypto jwe cripto module
type JWECrypto struct{}

// Encrypt Encrypts the content in a JWE Wrapper
func (c JWECrypto) Encrypt(k *rsa.PublicKey, m string) ([]byte, error) {
	msg, err := jwe.Encrypt([]byte(m), jwa.RSA_OAEP_256, k, jwa.A256CBC_HS512, jwa.NoCompress)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

// Decrypt Decrypts the JWE and return de message
func (c *JWECrypto) Decrypt(k *rsa.PrivateKey, m string) ([]byte, error) {
	msg, err := jwe.Decrypt([]byte(m), jwa.RSA_OAEP_256, k)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}
