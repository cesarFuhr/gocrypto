package keys

import (
	"reflect"
	"testing"
)

func FindKeyTest(t *testing.T) {
	keyRepo := InMemoryKeyRepository{}

	t.Run("Should return a Key", func(t *testing.T) {
		got, _ := keyRepo.FindKey("1")
		want := Key{}

		assertType(t, got, want)
	})
	t.Run("Should return the right key", func(t *testing.T) {
		iniKey := Key{ID: "1"}
		iniKeyRepo := InMemoryKeyRepository{map[string]Key{"1": iniKey}}

		got, _ := iniKeyRepo.FindKey("1")
		want := iniKey

		assertType(t, got, want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func InsertKeyTest(t *testing.T) {
	keyRepo := InMemoryKeyRepository{}

	t.Run("Should insert a Key", func(t *testing.T) {
		key := Key{ID: "1"}
		keyRepo.InsertKey(key)

		got, _ := keyRepo.FindKey(key.ID)
		want := key

		assertType(t, got, want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
