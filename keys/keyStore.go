package keys

import (
	"crypto/rsa"
	"time"
)

type Keys struct {
	Scope      string
	Expiration time.Time
	Priv       *rsa.PrivateKey
	Pub        *rsa.PublicKey
}

type KeySource interface {
	Take() *rsa.PrivateKey
}

type KeyStore struct {
	source KeySource
}

func (s *KeyStore) CreateKeys(scope string, expiration time.Time) Keys {
	keys := Keys{}
	keys.Priv = s.source.Take()
	keys.Pub = &keys.Priv.PublicKey
	keys.Scope = scope
	keys.Expiration = expiration
	return keys
}
