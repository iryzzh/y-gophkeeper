package sqlite

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/store"
)

func makeUser(tb testing.TB, id ...string) *models.User {
	tb.Helper()

	var uid string
	if len(id) > 0 {
		uid = id[0]
	} else {
		uid = uuid.New().String()
	}

	return &models.User{
		ID:           uid,
		Login:        "test-user",
		Password:     "test-password",
		PasswordHash: "test-password-hash",
	}
}

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.User
		wantErr error
	}{
		{
			name:    "create ok",
			want:    makeUser(t),
			wantErr: nil,
		},
		{
			name:    "already exist",
			want:    makeUser(t),
			wantErr: store.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				db: setupStore(t),
			}
			defer func() { _ = r.db.Close() }()

			if tt.wantErr == store.ErrUserAlreadyExists {
				_ = r.Create(context.Background(), tt.want)
			}

			err := r.Create(context.Background(), tt.want)
			if err != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.User
		wantErr error
	}{
		{
			name:    "found",
			want:    makeUser(t),
			wantErr: nil,
		},
		{
			name:    "not found",
			want:    &models.User{},
			wantErr: store.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				db: setupStore(t),
			}
			defer func() { _ = r.db.Close() }()

			if tt.wantErr != store.ErrUserNotFound {
				err := r.Create(context.Background(), tt.want)
				if err != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			got, err := r.FindByID(context.Background(), tt.want.ID)
			if err != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want.Password = ""

			if got == nil {
				got = &models.User{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_FindByLogin(t *testing.T) {
	tests := []struct {
		name    string
		want    *models.User
		wantErr error
	}{
		{
			name:    "found",
			want:    makeUser(t),
			wantErr: nil,
		},
		{
			name:    "not found",
			want:    &models.User{},
			wantErr: store.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				db: setupStore(t),
			}
			defer func() { _ = r.db.Close() }()

			if tt.wantErr != store.ErrUserNotFound {
				err := r.Create(context.Background(), tt.want)
				if err != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			got, err := r.FindByLogin(context.Background(), tt.want.Login)
			if err != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want.Password = ""

			if got == nil {
				got = &models.User{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
