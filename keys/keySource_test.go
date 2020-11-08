package keys

import (
	"crypto/rsa"
	"io"
	"sync"
	"testing"
)

type keyGeneratorStub struct {
	called int
	stored int
	wg     sync.WaitGroup
	mu     sync.Mutex
}

func (g *keyGeneratorStub) GenerateKey(r io.Reader, s int) (*rsa.PrivateKey, error) {
	g.mu.Lock()
	g.called++
	g.mu.Unlock()

	return rsa.GenerateKey(r, s)
}

func TestSyncTake(t *testing.T) {
	keySource := SynchronousKeySource{}

	t.Run("Should return a valid rsa PrivKey", func(t *testing.T) {
		got, _ := keySource.Take()
		want := mockKeys

		assertType(t, got, want)
	})
}

func TestPoolTake(t *testing.T) {
	var waitGroup sync.WaitGroup
	keyGenStub := keyGeneratorStub{wg: waitGroup}
	keySource := PoolKeySource{make(chan *rsa.PrivateKey, 2), &keyGenStub}

	t.Run("returns a valid rsa PrivKey", func(t *testing.T) {
		got, _ := keySource.Take()
		want := mockKeys

		<-keySource.pool
		assertType(t, got, want)
	})
	t.Run("if there is no keys in the pool calls GenerateKey, and create one ascyncronouslly", func(t *testing.T) {
		keyGenStub.called = 0
		got, _ := keySource.Take()
		want := mockKeys

		<-keySource.pool
		assertType(t, got, want)
	})
	t.Run("if there is keys in the pool should not call GenerateKey, pop one key and create one ascyncronouslly", func(t *testing.T) {
		keySource.pool <- mockKeys
		keyGenStub.called = 0
		keySource.Take()

		<-keySource.pool
		assertValue(t, keyGenStub.called, 1)
		assertValue(t, len(keySource.pool), 0)
	})
}

func assertValue(t *testing.T, g, w interface{}) {
	t.Helper()

	if w != g {
		t.Errorf("want %v, got %v", w, g)
	}
}
