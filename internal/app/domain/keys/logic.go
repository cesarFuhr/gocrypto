package keys

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// KeyService Stores keys giving scopes and type
type KeyService struct {
	Source KeySource
	Repo   KeyRepository
}

// NewKeyService creates a new KeyService
func NewKeyService(s KeySource, r KeyRepository) *KeyService {
	return &KeyService{
		Source: s,
		Repo:   r,
	}
}

// CreateKey Creates a Key, scoping it and setting the expiration
func (s *KeyService) CreateKey(scope string, expiration time.Time) (Key, error) {
	newKey, _ := s.Source.Take()
	key := Key{
		Priv:       newKey,
		Pub:        &newKey.PublicKey,
		Scope:      scope,
		Expiration: expiration,
		ID:         uuid.New().String(),
	}

	if err := s.Repo.InsertKey(key); err != nil {
		return Key{}, err
	}

	return key, nil
}

var (
	// ErrKeyNotFound the Key with the requested ID was not found in this store
	ErrKeyNotFound = errors.New("requested key was not found")
	// ErrKeyOutOfScope the Key was found but is not within the requested scope
	ErrKeyOutOfScope = errors.New("requested key is out of scope")
)

// FindKey Finds a key by ID
func (s *KeyService) FindKey(keyID string) (Key, error) {
	key, err := s.Repo.FindKey(keyID)
	if err != nil {
		if err == ErrKeyNotFound {
			return Key{}, ErrKeyNotFound
		}
		return Key{}, err
	}
	return key, nil
}

// FindScopedKey Find a key by ID within the scope
func (s *KeyService) FindScopedKey(keyID string, scope string) (Key, error) {
	key, err := s.FindKey(keyID)
	if err != nil {
		return Key{}, err
	}
	if key.Scope != scope {
		return Key{}, ErrKeyOutOfScope
	}

	return key, nil
}

// FindKeysByScope Find a key by ID within the scope
func (s *KeyService) FindKeysByScope(scope string) ([]Key, error) {
	keys, err := s.Repo.FindKeysByScope(scope)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
