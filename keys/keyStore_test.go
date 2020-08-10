package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"
	"time"
)

type KeySourceStub struct {
}

func (p KeySourceStub) Take() (privKey *rsa.PrivateKey) {
	privKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	return
}

func TestCreateKeys(t *testing.T) {
	keyStore := KeyStore{
		source: KeySourceStub{},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got := keyStore.CreateKeys("scope", time.Now())
		want := Keys{}
		want.Priv, _ = rsa.GenerateKey(rand.Reader, 2048)
		want.Pub = &want.Priv.PublicKey

		assertType(t, got, want)
		assertType(t, got.Priv, want.Priv)
		assertType(t, got.Pub, want.Pub)
	})
	t.Run("Should return expiration date", func(t *testing.T) {
		key := keyStore.CreateKeys("scope", time.Now().AddDate(0, 0, 1))
		got := key.Expiration

		assertTime(t, got, time.Now().AddDate(0, 0, 1))
	})
	t.Run("returned Keys should have the scope property", func(t *testing.T) {
		keys := keyStore.CreateKeys("scope", time.Now())
		got := keys.Scope
		want := "scope"

		assertValue(t, got, want)
	})
}

func assertType(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("got %T want %T", got, want)
	}
}

func assertValue(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertTime(t *testing.T, got, want time.Time) {
	t.Helper()
	if got.Round(time.Second) != want.Round(time.Second) {
		t.Errorf("got %v want %v", got, want)
	}
}
