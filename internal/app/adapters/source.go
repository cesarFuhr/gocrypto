package adapters

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

// KeyGenerator Type created to be a proxy of rsa key generator
type KeyGenerator struct{}

// GenerateKey proxy generates RSA Keys
func (g *KeyGenerator) GenerateKey(r io.Reader, s int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(r, s)
}

// NewPoolKeySource creates initialized PoolKeySource
func NewPoolKeySource(keySize int, keyPoolSize int) PoolKeySource {
	var s PoolKeySource
	s.keySize = keySize
	s.Kgen = &KeyGenerator{}
	s.Pool = make(chan *rsa.PrivateKey, keyPoolSize)
	return s
}

// PoolKeySource a key source based on a pool of keys
type PoolKeySource struct {
	Pool    chan *rsa.PrivateKey
	Kgen    keyGenerator
	keySize int
}

// Take Takes one key from the source
func (s *PoolKeySource) Take() (*rsa.PrivateKey, error) {
	defer func() { go s.addKeyToPoll() }()
	if len(s.Pool) > 0 {
		k := <-s.Pool
		return k, nil
	}
	return s.Kgen.GenerateKey(rand.Reader, s.keySize)
}

func (s *PoolKeySource) addKeyToPoll() {
	if len(s.Pool) < cap(s.Pool) {
		k, _ := s.Kgen.GenerateKey(rand.Reader, s.keySize)
		s.Pool <- k
	}
}

// WarmUp fills the bufered channel with keys
func (s *PoolKeySource) WarmUp() {
	for len(s.Pool) < cap(s.Pool) {
		k, err := s.Kgen.GenerateKey(rand.Reader, s.keySize)
		if err != nil {
			panic(err)
		}
		s.Pool <- k
	}
}
