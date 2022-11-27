package token

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/store"
	"github.com/iryzzh/gophkeeper/internal/store/sqlite"
	"github.com/iryzzh/gophkeeper/internal/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func testStore(t *testing.T) store.Store {
	t.Helper()
	cfg, err := utils.TestConfig(t)
	if err != nil {
		t.Fatal(err)
	}
	st, err := sqlite.NewStore(cfg.DB.DSN, cfg.DB.MigrationsPath)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestService_CreateToken(t *testing.T) {
	type fields struct {
		accessSecret  []byte
		refreshSecret []byte
		atExpiresIn   int
		rtExpiresIn   int
	}
	tests := []struct {
		name    string
		fields  fields
		user    *models.User
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				accessSecret:  []byte("Access-Super-Secret"),
				refreshSecret: []byte("Refresh-Super-Secret"),
				atExpiresIn:   15,
				rtExpiresIn:   10080,
			},
			user: &models.User{
				ID:           uuid.NewString(),
				Login:        "test-user",
				PasswordHash: "test-password",
			},
			wantErr: nil,
		},
		{
			name:    "invalid user",
			user:    nil,
			wantErr: ErrInvalidUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				store:         testStore(t),
				accessSecret:  tt.fields.accessSecret,
				refreshSecret: tt.fields.refreshSecret,
				atExpiresIn:   tt.fields.atExpiresIn,
				rtExpiresIn:   tt.fields.rtExpiresIn,
			}
			defer func() { _ = s.store.Close() }()
			_, err := s.Create(context.Background(), tt.user)
			require.Equalf(t, tt.wantErr, err, "create token expected error = %v, got error = %v", tt.wantErr, err)
		})
	}
}

func TestService_ValidateToken(t *testing.T) {
	type fields struct {
		accessSecret  []byte
		refreshSecret []byte
		atExpiresIn   int
		rtExpiresIn   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
		user    func(st store.Store) *models.User
	}{
		{
			name: "ok",
			fields: fields{
				accessSecret:  []byte("Access-Super-Secret"),
				refreshSecret: []byte("Refresh-Super-Secret"),
				atExpiresIn:   15,
				rtExpiresIn:   10080,
			},
			user: func(st store.Store) *models.User {
				user, err := st.User().Create(context.Background(), &models.User{
					ID:           uuid.NewString(),
					Login:        "test-user",
					PasswordHash: "test-password",
				})
				if err != nil {
					t.Fatal(err)
				}

				return user
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				store:         testStore(t),
				accessSecret:  tt.fields.accessSecret,
				refreshSecret: tt.fields.refreshSecret,
				atExpiresIn:   tt.fields.atExpiresIn,
				rtExpiresIn:   tt.fields.rtExpiresIn,
			}
			defer func() { _ = s.store.Close() }()
			user := tt.user(s.store)
			token := &models.Token{}
			var err error
			if user != nil {
				token, err = s.Create(context.Background(), user)
				require.NoError(t, err)
			}
			_, err = s.Validate(context.Background(), token.AccessToken)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_RefreshToken(t *testing.T) {
	type fields struct {
		accessSecret  []byte
		refreshSecret []byte
		atExpiresIn   int
		rtExpiresIn   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
		user    func(st store.Store) *models.User
	}{
		{
			name: "ok",
			fields: fields{
				accessSecret:  []byte("Access-Super-Secret"),
				refreshSecret: []byte("Refresh-Super-Secret"),
				atExpiresIn:   15,
				rtExpiresIn:   10080,
			},
			user: func(st store.Store) *models.User {
				user, err := st.User().Create(context.Background(), &models.User{
					ID:           uuid.NewString(),
					Login:        "test-user",
					PasswordHash: "test-password",
				})
				if err != nil {
					t.Fatal(err)
				}

				return user
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				store:         testStore(t),
				accessSecret:  tt.fields.accessSecret,
				refreshSecret: tt.fields.refreshSecret,
				atExpiresIn:   tt.fields.atExpiresIn,
				rtExpiresIn:   tt.fields.rtExpiresIn,
			}
			defer func() { _ = s.store.Close() }()
			user := tt.user(s.store)
			token := &models.Token{}
			var err error
			if user != nil {
				token, err = s.Create(context.Background(), user)
				require.NoError(t, err)
			}
			s.atExpiresIn = 1
			newToken, rtErr := s.Refresh(context.Background(), token.RefreshToken)
			require.Equal(t, tt.wantErr, rtErr)
			require.NotEqual(t, token, newToken)
		})
	}
}

func TestNewTokenService(t *testing.T) {
	s := NewService(testStore(t), 0, 0, []byte("a"), []byte("a"))
	require.NotNil(t, s)
}

func Test_parseJWT(t *testing.T) {
	type args struct {
		tokenStr string
		secret   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *jwt.Token
		wantErr error
	}{
		{
			name:    "invalid",
			wantErr: ErrInvalidToken,
		},
		{
			name: "ok",
			args: args{
				tokenStr: func() string {
					claims := &jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
						Issuer:    "test",
					}

					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					ss, err := token.SignedString([]byte("AllYourBase"))
					if err != nil {
						t.Fatal(err)
					}
					return ss
				}(),
				secret: []byte("AllYourBase"),
			},
			wantErr: ErrTokenExpired,
		},
		{
			name: "bad signing method",
			args: args{
				tokenStr: func() string {
					claims := &jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					}
					_, key, _ := ed25519.GenerateKey(rand.Reader)
					token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
					ss, err := token.SignedString(key)
					if err != nil {
						t.Fatal(err)
					}
					return ss
				}(),
				secret: []byte("AllYourBase"),
			},
			wantErr: ErrTokenBadSigningMethod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseJWT(tt.args.tokenStr, tt.args.secret)
			if !errors.Is(err, tt.wantErr) {
				t.Logf("want err = %v, got = %v", tt.wantErr, err.Error())
			}
		})
	}
}
