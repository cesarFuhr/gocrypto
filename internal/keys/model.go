package keys

import (
	"crypto/rsa"
	"time"
)

// Key Representation of a rsa key with scope, ID and expiration
type Key struct {
	Scope      string
	ID         string
	Expiration time.Time
	Priv       *rsa.PrivateKey
	Pub        *rsa.PublicKey
}
