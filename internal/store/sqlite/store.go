package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iryzzh/gophkeeper/internal/store"
	"github.com/pkg/errors"
)

// Store is a store.
type Store struct {
	db *sql.DB
}

func NewStore(dsn, migrationsPath string) (*Store, error) {
	db, _ := sql.Open("sqlite3", dsn)
	s := &Store{
		db: db,
	}

	if err := s.migrate(migrationsPath); err != nil {
		return nil, err
	}

	return s, db.Ping()
}

// Close closes the database and prevents new queries from starting.
func (s *Store) Close() error {
	return s.db.Close()
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) migrate(path string) error {
	f := path
	if f == "" {
		return nil
	}
	_, err := os.Stat(f)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("migration file not found: %v", f)
		}
		return err
	}

	driver, err := sqlite3.WithInstance(s.db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	if m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", f),
		"sqlite3", driver); err == nil {
		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	return err
}

func (s *Store) User() store.UserRepository {
	return &UserRepository{db: s.db}
}

func (s *Store) Item() store.ItemRepository {
	return &ItemRepository{db: s.db}
}
