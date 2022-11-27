package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/store"
)

var (
	// ErrPasswordCannotBeEmpty returns when the password is empty.
	ErrPasswordCannotBeEmpty = errors.New("password cannot be empty")
	// ErrInvalidUser returns when the user is invalid.
	ErrInvalidUser = errors.New("invalid user")
	// ErrLoginOrPasswordIsInvalid returns when the login or password is invalid.
	ErrLoginOrPasswordIsInvalid = errors.New("login or password is invalid")
	// ErrUserNotFound returns when the user is not found.
	ErrUserNotFound = errors.New("user not found")
)

// Service is a service for user interaction.
type Service struct {
	store          store.Store
	argon2idParams *argon2id.Params
}

// NewService creates a new service.
func NewService(s store.Store, hashMemory, hashIterations uint32, hashParallelism uint8, saltLength, keyLength uint32) *Service {
	return &Service{
		store: s,
		argon2idParams: &argon2id.Params{
			Memory:      hashMemory,
			Iterations:  hashIterations,
			Parallelism: hashParallelism,
			SaltLength:  saltLength,
			KeyLength:   keyLength,
		},
	}
}

// Create encrypts the received password, creates a new user in the database and, if successful, returns the
// `models.User` struct
func (s *Service) Create(ctx context.Context, user, password string) (*models.User, int, error) {
	if user == "" {
		return nil, http.StatusBadRequest, ErrInvalidUser
	}
	if password == "" {
		return nil, http.StatusBadRequest, ErrPasswordCannotBeEmpty
	}

	hash, err := argon2id.CreateHash(password, s.argon2idParams)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	u := &models.User{
		ID:           uuid.NewString(),
		Login:        user,
		PasswordHash: hash,
	}
	_, err = s.store.User().Create(ctx, u)
	if errors.Is(err, store.ErrUserAlreadyExists) {
		return nil, http.StatusConflict, err
	}

	return u, http.StatusCreated, nil
}

func (s *Service) Login(ctx context.Context, user, password string) (*models.User, error) {
	if user == "" {
		return nil, ErrInvalidUser
	}

	if password == "" {
		return nil, ErrPasswordCannotBeEmpty
	}

	u, err := s.Find(ctx, user)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if ok, _ := argon2id.ComparePasswordAndHash(password, u.PasswordHash); !ok {
		return nil, ErrLoginOrPasswordIsInvalid
	}

	return u, nil
}

// Find looks for a user in the database by his login and returns the user if found.
func (s *Service) Find(ctx context.Context, user string) (*models.User, error) {
	return s.store.User().FindByLogin(ctx, user)
}
