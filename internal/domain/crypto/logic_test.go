package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/domain/keys"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

var (
	rsaKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	key       = keys.Key{
		Scope:      "scope",
		ID:         "id",
		Expiration: time.Now().AddDate(0, 0, 1),
		Priv:       rsaKey,
		Pub:        &rsaKey.PublicKey,
	}
)

type RepositoryStub struct{}

func (r *RepositoryStub) FindKey(id string) (keys.Key, error) {
	return key, nil
}

func TestCryptoEncrypt(t *testing.T) {
	crypto := NewCryptoService(&RepositoryStub{})
	t.Run("Should return a valid JWE", func(t *testing.T) {
		got, _ := crypto.Encrypt("id", "testingOK")

		if _, err := jwe.Decrypt(got, jwa.RSA_OAEP_256, key.Priv); err != nil {
			t.Errorf("Invalid jwe: %v", err)
		}
	})
	t.Run("Should be able to decrypt back", func(t *testing.T) {
		want := "test"
		encrypted, _ := crypto.Encrypt("id", want)

		decrypted, _ := jwe.Decrypt(encrypted, jwa.RSA_OAEP_256, key.Priv)
		got := string(decrypted)

		if want != got {
			t.Errorf("want %v, got %v", want, string(got))
		}
	})
}

func TestCryptoDecrypt(t *testing.T) {
	crypto := cryptoService{&RepositoryStub{}}
	t.Run("Should be able to decrypt a encrypted message", func(t *testing.T) {
		want := "test"
		encrypted, _ := crypto.Encrypt("id", want)

		decrypted, _ := crypto.Decrypt("id", string(encrypted))
		got := string(decrypted)

		if want != got {
			t.Errorf("want %v, got %v", want, string(got))
		}
	})
}
