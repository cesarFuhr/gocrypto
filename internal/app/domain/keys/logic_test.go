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
		return Key{}, ErrKeyNotFound
	}
	return key, nil
}

func (r *KeyRepositoryStub) FindKeysByScope(scope string) ([]Key, error) {
	if scope == "not found" {
		return nil, nil
	}

	return []Key{keyStub, keyStub}, nil
}

func (r *KeyRepositoryStub) InsertKey(key Key) error {
	r.store[key.ID] = key
	return nil
}

type KeySourceStub struct {
}

var (
	mockKeys, mockErr = rsa.GenerateKey(rand.Reader, 2048)
	keyStub           = Key{
		ID:         "id",
		Scope:      "scope",
		Expiration: time.Now(),
		Priv:       mockKeys,
		Pub:        &mockKeys.PublicKey,
	}
)

func (p *KeySourceStub) Take() (*rsa.PrivateKey, error) {
	return mockKeys, mockErr
}

func TestCreateKey(t *testing.T) {
	keyStore := KeyService{
		Source: &KeySourceStub{},
		Repo:   &KeyRepositoryStub{map[string]Key{}},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got, _ := keyStore.CreateKey("scope", time.Now())
		want := Key{}
		want.Priv, _ = rsa.GenerateKey(rand.Reader, 2048)
		want.Pub = &want.Priv.PublicKey

		assertType(t, got, want)
		assertType(t, got.Priv, want.Priv)
		assertType(t, got.Pub, want.Pub)
	})
	t.Run("Should return expiration date", func(t *testing.T) {
		key, _ := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))
		got := key.Expiration

		assertTime(t, got, time.Now().AddDate(0, 0, 1))
	})
	t.Run("returned Keys should have the scope property", func(t *testing.T) {
		key, _ := keyStore.CreateKey("scope", time.Now())
		got := key.Scope
		want := "scope"

		assertString(t, got, want)
	})
}

func TestFindKey(t *testing.T) {
	keyStore := KeyService{
		Source: &KeySourceStub{},
		Repo:   &KeyRepositoryStub{map[string]Key{}},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got, _ := keyStore.FindKey("id")
		want, _ := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))

		assertType(t, got, want)
	})
	t.Run("Should return the correct keypair", func(t *testing.T) {
		key, _ := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))
		found, _ := keyStore.FindKey(key.ID)

		assertString(t, found.ID, key.ID)
	})
	t.Run("Should return an error if key was not found", func(t *testing.T) {
		_, err := keyStore.FindKey("inexistent key.ID")

		if err != ErrKeyNotFound {
			t.Fatalf("was expecting a ErrKeyNotFound and didn't received")
		}
	})
}

func TestFindScopedKey(t *testing.T) {
	keyStore := KeyService{
		Source: &KeySourceStub{},
		Repo:   &KeyRepositoryStub{map[string]Key{}},
	}
	t.Run("Should return a keypair", func(t *testing.T) {
		got, _ := keyStore.FindScopedKey("id", "scope")
		want, _ := keyStore.CreateKey("scope", time.Now().AddDate(0, 0, 1))

		assertType(t, got, want)
	})
	t.Run("Should return an error if Key is out of scope", func(t *testing.T) {
		key, _ := keyStore.CreateKey("scope1", time.Now().AddDate(0, 0, 1))
		_, err := keyStore.FindScopedKey(key.ID, "scope2")

		if err != ErrKeyOutOfScope {
			t.Fatalf("was expecting a ErrKeyOutOfScope and received %v", err)
		}
	})
	t.Run("Should return an error if key was not found", func(t *testing.T) {
		_, err := keyStore.FindScopedKey("inexistent key.ID", "scope")

		if err != ErrKeyNotFound {
			t.Fatalf("was expecting a KeyNotFoundError and didn't received")
		}
	})
}

func TestFindKeysByScope(t *testing.T) {
	keyStore := KeyService{
		Source: &KeySourceStub{},
		Repo:   &KeyRepositoryStub{map[string]Key{}},
	}
	t.Run("Should return a slice of keypair", func(t *testing.T) {
		got, _ := keyStore.FindKeysByScope("scope")
		want := []Key{}

		assertType(t, got, want)
	})
	t.Run("Should not return an error if no key was not found", func(t *testing.T) {
		_, err := keyStore.FindKeysByScope("not found")

		if err != nil {
			t.Fatalf("was expecting a nil and didn't received")
		}
	})
}

func TestNewKeyService(t *testing.T) {
	source := &KeySourceStub{}
	repo := &KeyRepositoryStub{}
	t.Run("Returns a valid KeyService", func(t *testing.T) {
		got := NewKeyService(source, repo)
		want := &KeyService{}
		assertType(t, got, want)
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
