package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file" // fs source
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/iryzzh/y-gophkeeper/internal/store"
)

var (
	//go:embed migrations/*.sql
	fs embed.FS
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

	if err := s.migrate(dsn, migrationsPath); err != nil {
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

// IsUsersExist returns true if any users exist in the database.
func (s *Store) IsUsersExist() (result bool, err error) {
	err = s.db.QueryRow(`select exists (select 1 from users) AS result`).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

// IsItemsExist returns true if any items exist in the database.
func (s *Store) IsItemsExist() (result bool, err error) {
	err = s.db.QueryRow(`select exists (select 1 from items) AS result`).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

// migrate uses dsn and the migration file path as parameters. if
// the path does not exist, or some error occurs while trying to
// read the specified path, then the embedded filesystem with the
// default migration configuration is used.
func (s *Store) migrate(dsn, path string) error {
	var m *migrate.Migrate

	_, err := os.Stat(path)
	if err == nil {
		driver, _ := sqlite3.WithInstance(s.db, &sqlite3.Config{})
		m, _ = migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", path), "sqlite3", driver)
	} else {
		// log.Printf("error during an attempt to read the migration file path: %v. default migrations are used.", err)
		var d source.Driver
		d, _ = iofs.New(fs, "migrations")
		m, _ = migrate.NewWithSourceInstance("iofs", d, "sqlite3://"+dsn)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *Store) User() store.UserRepository {
	return &UserRepository{db: s.db}
}

func (s *Store) Item() store.ItemRepository {
	return &ItemRepository{db: s.db}
}
