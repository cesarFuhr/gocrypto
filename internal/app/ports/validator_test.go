package ports

import "testing"

var o = keyOpts{
	Scope:      "string",
	Expiration: "2021-11-20T18:01:24.100+00:00",
}

// func TestValidateWithTags(t *testing.T) {
// 	t.Run("returns nill if all properties are correct", func(t *testing.T) {
// 		got := validateWithTags(o)

// 		if got != nil {
// 			t.Errorf("got %v, want %v", got, nil)
// 		}
// 	})
// }

func TestValidateWithFunc(t *testing.T) {
	t.Run("returns nill if all properties are correct", func(t *testing.T) {
		got := validateWithFunc(o)

		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})
}

func TestValidateWithOzzo(t *testing.T) {
	t.Run("returns nill if all properties are correct", func(t *testing.T) {
		got := validateOzzo(o)

		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})
}

func TestValidateHome(t *testing.T) {
	t.Run("returns nill if all properties are correct", func(t *testing.T) {
		got := validateHomeSolution(o)

		if got != nil {
			t.Errorf("got %v, want %v", got, nil)
		}
	})
}

func BenchmarkWithTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validateWithTags(o)
	}
}

func BenchmarkHome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validateHomeSolution(o)
	}
}

func BenchmarkWithFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validateWithFunc(o)
	}
}

func BenchmarkOzzoSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validateOzzo(o)
	}
}
