package crypto

import "github.com/cesarFuhr/gocrypto/internal/app/domain/keys"

// Repository Persistency interface to serve cryptoService
type Repository interface {
	FindKey(string) (keys.Key, error)
}
