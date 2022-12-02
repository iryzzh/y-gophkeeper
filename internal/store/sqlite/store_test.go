package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iryzzh/y-gophkeeper/internal/config"
	"github.com/iryzzh/y-gophkeeper/internal/utils"
)

func setupStore(t *testing.T) *sql.DB {
	t.Helper()

	cfg, err := utils.TestConfig(t)
	if err != nil {
		t.Fatal(err)
	}

	st, err := NewStore(cfg.DB.DSN, cfg.DB.MigrationsPath)
	if err != nil {
		t.Fatal(err)
	}

	return st.db
}

func TestStore_Close(t *testing.T) {
	cfg, err := utils.TestConfig(t)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Close",
			fields: fields{
				db: func(cfg config.DBConfig) *sql.DB {
					db, err := sql.Open(cfg.Type, cfg.DSN)
					if err != nil {
						t.Fatal(err)
					}
					return db
				}(cfg.DB),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_User(t *testing.T) {
	cfg, err := utils.TestConfig(t)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		db  *sql.DB
		cfg config.DBConfig
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "User",
			fields: fields{
				db: func(cfg config.DBConfig) *sql.DB {
					db, err := sql.Open(cfg.Type, cfg.DSN)
					if err != nil {
						t.Fatal(err)
					}
					return db
				}(cfg.DB),
				cfg: cfg.DB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if got := s.User(); got == nil {
				t.Errorf("User() = %v", got)
			}
		})
	}
}

func TestStore_Ping(t *testing.T) {
	s := &Store{
		db: setupStore(t),
	}
	require.NoError(t, s.Ping())
}

func TestStore_Item(t *testing.T) {
	s := &Store{
		db: setupStore(t),
	}
	require.NotNil(t, s.Item())
}

func TestNewStore(t *testing.T) {
	tests := []struct {
		name           string
		dsn            string
		migrationsPath string
		wantErr        bool
	}{
		{
			name:           "migrate path is empty",
			dsn:            ":memory:",
			migrationsPath: "",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStore(tt.dsn, tt.migrationsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
