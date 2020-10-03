package keys

import (
	"crypto/rsa"
	"testing"
)

func InMemTakeTest(t *testing.T) {
	keySource := InMemoryKeySource{}

	t.Run("Should return a valid rsa PrivKey", func(t *testing.T) {
		got, _ := keySource.Take()
		want := rsa.PrivateKey{}

		assertType(t, got, want)
	})
}
