package keys

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Key Representation of a rsa key with scope, ID and expiration
type Key struct {
	Scope      string
	ID         string
	Expiration time.Time
	Priv       *rsa.PrivateKey
	Pub        *rsa.PublicKey
}

// KeySource Key provider of the KeyStore
type KeySource interface {
	Take() *rsa.PrivateKey
}

// KeyRepository Persistency interface to serve the KeyStore
type KeyRepository interface {
	FindKey(string) (Key, error)
	InsertKey(Key) error
}

// KeyStore Stores keys giving scopes and expiration
type KeyStore struct {
	source KeySource
	repo   KeyRepository
}

// CreateKey Creates a Key, scoping it and setting the expiration
func (s *KeyStore) CreateKey(scope string, expiration time.Time) Key {
	newKey := s.source.Take()
	key := Key{
		Priv:       newKey,
		Pub:        &newKey.PublicKey,
		Scope:      scope,
		Expiration: expiration,
		ID:         uuid.New().String(),
	}

	s.repo.InsertKey(key)

	return key
}

// ErrKeyNotFound the Key with the requested ID was not found in this store
var ErrKeyNotFound = errors.New("requested key was not found")

// ErrKeyOutOfScope the Key was found but is not within the requested scope
var ErrKeyOutOfScope = errors.New("requested key is out of scope")

// FindKey Finds a key by ID
func (s *KeyStore) FindKey(keyID string) (Key, error) {
	key, err := s.repo.FindKey(keyID)
	if err != nil {
		if err == ErrKeyNotFound {
			return Key{}, ErrKeyNotFound
		}
		return Key{}, err
	}
	return key, nil
}

// FindScopedKey Find a key by ID within the scope
func (s *KeyStore) FindScopedKey(keyID string, scope string) (Key, error) {
	key, err := s.FindKey(keyID)
	if err != nil {
		return Key{}, err
	}
	if key.Scope != scope {
		return Key{}, ErrKeyOutOfScope
	}

	return key, nil
}
