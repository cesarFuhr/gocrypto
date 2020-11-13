package keys

import (
	"crypto/x509"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestMemFindKey(t *testing.T) {
	keyRepo := InMemoryKeyRepository{map[string]Key{}}

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

func TestMemInsertKey(t *testing.T) {
	keyRepo := InMemoryKeyRepository{map[string]Key{}}

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

var key = Key{
	ID:         uuid.New().String(),
	Scope:      "scope",
	Expiration: time.Now().AddDate(0, 0, 1),
	Priv:       mockKeys,
	Pub:        &mockKeys.PublicKey,
}

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
func TestSQLInsertKey(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLKeyRepository{db: db}
	defer db.Close()

	t.Run("calls db.Exec with the right params", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO keys").WithArgs(
			key.ID,
			key.Scope,
			key.Expiration,
			anyTime{},
			x509.MarshalPKCS1PrivateKey(key.Priv),
			x509.MarshalPKCS1PublicKey(key.Pub),
		)

		repo.InsertKey(key)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectExec("INSERT INTO keys").WithArgs(
			key.ID,
			key.Scope,
			key.Expiration,
			anyTime{},
			x509.MarshalPKCS1PrivateKey(key.Priv),
			x509.MarshalPKCS1PublicKey(key.Pub),
		).WillReturnError(want)

		got := repo.InsertKey(key)

		assertValue(t, got, want)
	})
}

func TestSQLFindKey(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := SQLKeyRepository{db: db}
	defer db.Close()

	t.Run("calls db.QueryRow with the right params", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT id, scope, expiration, priv, pub
				FROM keys
				WHERE id`).WithArgs(key.ID)

		repo.FindKey(key.ID)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("SQL expectations failed: %s", err)
		}
	})

	t.Run("returns a complete Key object", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "scope", "expiration", "priv", "pub"}).
			AddRow(key.ID, key.Scope, key.Expiration,
				x509.MarshalPKCS1PrivateKey(key.Priv),
				x509.MarshalPKCS1PublicKey(key.Pub))
		mock.
			ExpectQuery(`
				SELECT id, scope, expiration, priv, pub
					FROM keys
					WHERE id`).
			WithArgs(key.ID).
			WillReturnRows(rows)

		returned, err := repo.FindKey(key.ID)

		assertValue(t, err, nil)
		if !reflect.DeepEqual(key, returned) {
			t.Errorf("want %v, got %v", key, returned)
		}
	})

	t.Run("proxys the error from the sql db", func(t *testing.T) {
		want := errors.New("an error")
		mock.ExpectQuery(`
			SELECT id, scope, expiration, priv, pub
				FROM keys
				WHERE id`).WithArgs(key.ID).WillReturnError(want)

		_, got := repo.FindKey(key.ID)

		assertValue(t, got, want)
	})

	t.Run("not founding the key, return a ErrKeyNotFound", func(t *testing.T) {
		want := ErrKeyNotFound
		mock.ExpectQuery(`
			SELECT id, scope, expiration, priv, pub
				FROM keys
				WHERE id`).WithArgs(key.ID).WillReturnRows(sqlmock.NewRows([]string{}))

		_, got := repo.FindKey(key.ID)

		assertValue(t, got, want)
	})
}
