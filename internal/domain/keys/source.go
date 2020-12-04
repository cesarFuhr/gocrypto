package keys

import "crypto/rsa"

// KeySource Key provider of the KeyStore
type KeySource interface {
	Take() (*rsa.PrivateKey, error)
}
