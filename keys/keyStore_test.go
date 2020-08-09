package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"
)

type KeySourceStub struct {
}

func (p KeySourceStub) Pop() (privKey *rsa.PrivateKey) {
	privKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	return
}

func TestCreateKeys(t *testing.T) {
	keyStore := KeyStore{
		source: KeySourceStub{},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got := keyStore.CreateKeys("scope")
		want := Keys{}
		want.Priv, _ = rsa.GenerateKey(rand.Reader, 2048)
		want.Pub = &want.Priv.PublicKey

		assertType(t, got, want)
		assertType(t, got.Priv, want.Priv)
		assertType(t, got.Pub, want.Pub)
	})
	t.Run("returned Keys should have the scope information", func(t *testing.T) {
		keys := keyStore.CreateKeys("scope")
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
