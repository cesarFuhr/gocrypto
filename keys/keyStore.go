package keys

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Key struct {
	Scope      string
	ID         string
	Expiration time.Time
	Priv       *rsa.PrivateKey
	Pub        *rsa.PublicKey
}

type KeySource interface {
	Take() *rsa.PrivateKey
}

type KeyRepository interface {
	FindKey(string) (Key, error)
	InsertKey(Key) error
}

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

// FindScopedKey Find a key by ID within the scope
func (s *KeyStore) FindScopedKey(keyID string, scope string) (Key, error) {
	key, err := s.repo.FindKey(keyID)
	if err != nil {
		if err == KeyNotFoundError {
			return Key{}, KeyNotFoundError
		}
		return Key{}, err
	}
	return key, nil
}

var KeyNotFoundError = errors.New("requested key was not found")
