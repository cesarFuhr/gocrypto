package crypto

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

type CryptoService struct {
	repo Repository
}

// NewCryptoService creates a new crypto service
func NewCryptoService(r Repository) CryptoService {
	return CryptoService{
		repo: r,
	}
}

// Encrypt Encrypts the content in a JWE Wrapper
func (s *CryptoService) Encrypt(keyID string, m string) ([]byte, error) {
	key, err := s.repo.FindKey(keyID)
	if err != nil {
		return []byte{}, err
	}

	msg, err := jwe.Encrypt([]byte(m), jwa.RSA_OAEP_256, key.Pub, jwa.A256CBC_HS512, jwa.NoCompress)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

// Decrypt Decrypts the JWE and return de message
func (s *CryptoService) Decrypt(keyID string, m string) ([]byte, error) {
	key, err := s.repo.FindKey(keyID)
	if err != nil {
		return []byte{}, err
	}

	msg, err := jwe.Decrypt([]byte(m), jwa.RSA_OAEP_256, key.Priv)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}
