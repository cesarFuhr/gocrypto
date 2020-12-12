package adapters

import (
	"crypto/x509"
	"database/sql"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/app/domain/keys"
	// Is there any other way?
	_ "github.com/lib/pq"
)

// InMemoryKeyRepository simple in memory key repository
type InMemoryKeyRepository struct {
	Store map[string]keys.Key
}

// FindKey finds and returns the requested key
func (r *InMemoryKeyRepository) FindKey(id string) (keys.Key, error) {
	key, ok := r.Store[id]
	if ok == false {
		return keys.Key{}, keys.ErrKeyNotFound
	}
	return key, nil
}

// InsertKey Inserts a key into the repository
func (r *InMemoryKeyRepository) InsertKey(key keys.Key) error {
	r.Store[key.ID] = key
	return nil
}

// NewSQLKeyRepository returns a new sql repository instance
func NewSQLKeyRepository(db *sql.DB) SQLKeyRepository {
	return SQLKeyRepository{db: db}
}

// SQLKeyRepository sql database persistency
type SQLKeyRepository struct {
	db *sql.DB
}

var findKeyStatement = `
	SELECT id, scope, expiration, priv, pub 
		FROM keys 
		WHERE id = $1`

// FindKey finds and returns the requested key
func (r *SQLKeyRepository) FindKey(id string) (keys.Key, error) {
	row := r.db.QueryRow(findKeyStatement, id)

	var k keys.Key
	var priv, pub []byte

	switch err := row.Scan(&k.ID, &k.Scope, &k.Expiration, &priv, &pub); err {
	case nil:
		k.Priv, err = x509.ParsePKCS1PrivateKey(priv)
		k.Pub, err = x509.ParsePKCS1PublicKey(pub)
		if err != nil {
			return keys.Key{}, err
		}
	case sql.ErrNoRows:
		return keys.Key{}, keys.ErrKeyNotFound
	default:
		return keys.Key{}, err
	}
	return k, nil
}

var insertKeyStatement = `
	INSERT INTO keys (id, scope, expiration, creation, priv, pub)
		VALUES ($1, $2, $3, $4, $5, $6)`

// InsertKey Inserts a key into the repository
func (r *SQLKeyRepository) InsertKey(k keys.Key) error {
	_, err := r.db.Exec(
		insertKeyStatement,
		k.ID,
		k.Scope,
		k.Expiration,
		time.Now(),
		x509.MarshalPKCS1PrivateKey(k.Priv),
		x509.MarshalPKCS1PublicKey(k.Pub),
	)
	return err
}
