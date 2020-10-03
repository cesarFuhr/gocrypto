package keys

import (
	"crypto/rand"
	"crypto/rsa"
)

type InMemoryKeySource struct {}

func (s *InMemoryKeySource) Take() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
