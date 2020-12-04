package keys

// KeyRepository Persistency interface to serve the KeyStore
type KeyRepository interface {
	FindKey(string) (Key, error)
	InsertKey(Key) error
}
