package keys

// InMemoryKeyRepository simple in memory key repository
type InMemoryKeyRepository struct {
	Store map[string]Key
}

// FindKey finds and returns the requested key
func (r *InMemoryKeyRepository) FindKey(id string) (Key, error) {
	key, ok := r.Store[id]
	if ok == false {
		return Key{}, ErrKeyNotFound
	}
	return key, nil
}

// InsertKey Inserts a key into the repository
func (r *InMemoryKeyRepository) InsertKey(key Key) error {
	r.Store[key.ID] = key
	return nil
}
