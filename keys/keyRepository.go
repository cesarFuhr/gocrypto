package keys

type InMemoryKeyRepository struct {
	Store map[string]Key
}

func (r *InMemoryKeyRepository) FindKey(id string) (Key, error) {
	key, ok := r.Store[id]
	if ok == false {
		return Key{}, ErrKeyNotFound
	}
	return key, nil
}

func (r *InMemoryKeyRepository) InsertKey(key Key) error {
	r.Store[key.ID] = key
	return nil
}
