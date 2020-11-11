package keys

import (
	"database/sql"
	"fmt"

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
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

// SQLKeyRepository sql database persistency
type SQLKeyRepository struct {
	Store map[string]Key
	Cfg   SQLConfigs
	db    *sql.DB
}

// FindKey finds and returns the requested key
func (r *SQLKeyRepository) FindKey(id string) (Key, error) {
	key, ok := r.Store[id]
	if ok == false {
		return Key{}, ErrKeyNotFound
	}
	return key, nil
}

// InsertKey Inserts a key into the repository
func (r *SQLKeyRepository) InsertKey(key Key) error {
	r.Store[key.ID] = key
	return nil
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
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	r.db = db
	fmt.Println("Connected to SQLKeyRepository")
	return nil
}
