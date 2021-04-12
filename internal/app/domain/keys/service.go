package keys

import "time"

// KeyService describes the keys service methods
type KeyService interface {
	CreateKey(string, time.Time) (Key, error)
	FindKey(string) (Key, error)
	FindScopedKey(string, string) (Key, error)
	FindKeysByScope(string) ([]Key, error)
}
