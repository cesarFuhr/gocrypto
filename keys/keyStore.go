package keys

import (
	"crypto/rsa"
)

type Keys struct {
	Scope string
	Priv  *rsa.PrivateKey
	Pub   *rsa.PublicKey
}

type KeyPool interface {
	Pop() *rsa.PrivateKey
}

type KeyStore struct {
	source KeyPool
}

func (s *KeyStore) CreateKeys(scope string) Keys {
	keys := Keys{}
	keys.Priv = s.source.Pop()
	keys.Pub = &keys.Priv.PublicKey
	keys.Scope = scope
	return keys
}
