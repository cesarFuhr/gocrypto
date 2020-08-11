package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"
	"time"
)

type KeyRepositoryStub struct {
	store map[string]Key
}

func (r *KeyRepositoryStub) FindKey(keyID string) (Key, error) {
	key, ok := r.store[keyID]
	if ok == false {
		return Key{}, KeyNotFoundError
	}
	return key, nil
}

func (r *KeyRepositoryStub) InsertKey(key Key) error {
	r.store[key.ID] = key
	return nil
}

type KeySourceStub struct {
}

var mockKeys, _ = rsa.GenerateKey(rand.Reader, 2048)

func (p KeySourceStub) Take() (privKey *rsa.PrivateKey) {
	privKey = mockKeys
	return
}

func TestCreateKey(t *testing.T) {
	keyStore := KeyStore{
		source: KeySourceStub{},
		repo:   &KeyRepositoryStub{map[string]Key{}},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got := keyStore.CreateKey("scope", time.Now())
		want := Key{}
		want.Priv, _ = rsa.GenerateKey(rand.Reader, 2048)
		want.Pub = &want.Priv.PublicKey

		assertType(t, got, want)
		assertType(t, got.Priv, want.Priv)
		assertType(t, got.Pub, want.Pub)
	})
	t.Run("Should return expiration date", func(t *testing.T) {
		key := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))
		got := key.Expiration

		assertTime(t, got, time.Now().AddDate(0, 0, 1))
	})
	t.Run("returned Keys should have the scope property", func(t *testing.T) {
		key := keyStore.CreateKey("scope", time.Now())
		got := key.Scope
		want := "scope"

		assertString(t, got, want)
	})
}

func TestFindScopedKeys(t *testing.T) {
	keyStore := KeyStore{
		source: KeySourceStub{},
		repo:   &KeyRepositoryStub{map[string]Key{}},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got, _ := keyStore.FindScopedKey("id", "scope")
		want := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))

		assertType(t, got, want)
	})
	t.Run("Should return the correct keypair", func(t *testing.T) {
		key := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))
		found, _ := keyStore.FindScopedKey(key.ID, "scope")

		assertString(t, found.ID, key.ID)
	})
	t.Run("Should return an error if key was not found", func(t *testing.T) {
		_, err := keyStore.FindScopedKey("inexistent key.ID", "scope")

		if err != KeyNotFoundError {
			t.Fatalf("was expecting a KeyNotFoundError and didn't received")
		}
	})
}

func assertType(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("got %T want %T", got, want)
	}
}

func assertString(t *testing.T, got, want string) {
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
