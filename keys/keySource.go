package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"io"
)

// SynchronousKeySource simple key source
type SynchronousKeySource struct{}

// Take Takes one key from the source
func (s *SynchronousKeySource) Take() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

type keyGenerator interface {
	GenerateKey(io.Reader, int) (*rsa.PrivateKey, error)
}

// PoolKeySource a key source based on a pool of keys
type PoolKeySource struct {
	pool chan *rsa.PrivateKey
	kgen keyGenerator
}

// Take Takes one key from the source
func (s *PoolKeySource) Take() (*rsa.PrivateKey, error) {
	go s.addKeyToPoll()
	if len(s.pool) > 0 {
		k := <-s.pool
		return k, nil
	}
	return s.kgen.GenerateKey(rand.Reader, 2048)
}

func (s *PoolKeySource) addKeyToPoll() {
	if len(s.pool) < cap(s.pool) {
		k, _ := s.kgen.GenerateKey(rand.Reader, 2048)
		s.pool <- k
	}
}
