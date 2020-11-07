package keys

import (
	"crypto/rand"
	"crypto/rsa"
)

// SynchronousKeySource simple key source
type SynchronousKeySource struct{}

// Take Takes one key from the source
func (s *SynchronousKeySource) Take() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
