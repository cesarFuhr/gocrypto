package keys

import (
	"crypto/x509"
	"database/sql"
	"fmt"
	"time"

	// Is there any other way?
	_ "github.com/lib/pq"
)

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

// SQLConfigs configuration for a sql database
type SQLConfigs struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	Dbname       string
	MaxOpenConns int
}

// SQLKeyRepository sql database persistency
type SQLKeyRepository struct {
	Cfg SQLConfigs
	db  *sql.DB
}

var findKeyStatement = `
	SELECT id, scope, expiration, priv, pub 
		FROM keys 
		WHERE id = $1`

// FindKey finds and returns the requested key
func (r *SQLKeyRepository) FindKey(id string) (Key, error) {
	row := r.db.QueryRow(findKeyStatement, id)

	var k Key
	var priv, pub []byte

	switch err := row.Scan(&k.ID, &k.Scope, &k.Expiration, &priv, &pub); err {
	case nil:
		k.Priv, err = x509.ParsePKCS1PrivateKey(priv)
		k.Pub, err = x509.ParsePKCS1PublicKey(pub)
		if err != nil {
			return Key{}, err
		}
	case sql.ErrNoRows:
		return Key{}, ErrKeyNotFound
	default:
		return Key{}, err
	}
	return k, nil
}

var insertKeyStatement = `
	INSERT INTO keys (id, scope, expiration, creation, priv, pub)
		VALUES ($1, $2, $3, $4, $5, $6)`

// InsertKey Inserts a key into the repository
func (r *SQLKeyRepository) InsertKey(k Key) error {
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

// Connect Created a connection with the database
func (r *SQLKeyRepository) Connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		r.Cfg.Host, r.Cfg.Port, r.Cfg.User, r.Cfg.Password, r.Cfg.Dbname)

	db, err := sql.Open(r.Cfg.Driver, psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(r.Cfg.MaxOpenConns)

	r.db = db
	fmt.Println("Connected to SQLKeyRepository")
	return nil
}
