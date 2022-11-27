package token

import (
	"crypto"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/store"
	"golang.org/x/net/context"
)

var (
	// ErrTokenBadSigningMethod returns when the token signing method does not match the expected.
	ErrTokenBadSigningMethod = errors.New("bad signing method")
	// ErrTokenExpired returns when the access or refresh token is expired.
	ErrTokenExpired = errors.New("token has expired")
	// ErrInvalidToken returns when the token is invalid.
	ErrInvalidToken = errors.New("token is invalid")
	// ErrInvalidUser returns when the user is invalid.
	ErrInvalidUser = errors.New("invalid user")
)

// Service is a token service that creates, validates and updates tokens, storing them in a database
type Service struct {
	store         store.Store
	accessSecret  []byte
	refreshSecret []byte
	atExpiresIn   int
	rtExpiresIn   int
}

// NewService returns a new token service.
func NewService(s store.Store, atExpiresIn, rtExpiresIn int, accessSecret, refreshSecret []byte) *Service {
	return &Service{
		store:         s,
		atExpiresIn:   atExpiresIn,
		rtExpiresIn:   rtExpiresIn,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

type claims struct {
	jwt.RegisteredClaims
	Login  string `json:"login"`
	UserID string `json:"user_id"`
}

// Create creates a new token.
func (s *Service) Create(_ context.Context, user *models.User) (*models.Token, error) {
	if user == nil {
		return nil, ErrInvalidUser
	}

	now := time.Now()

	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(s.atExpiresIn))),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Login:  user.Login,
		UserID: user.ID,
	}).SignedString(s.accessSecret)
	if err != nil {
		return nil, err
	}

	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(s.atExpiresIn))),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Login:  user.Login,
		UserID: user.ID,
	}).SignedString(s.refreshSecret)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		AccessToken:  at,
		RefreshToken: rt,
		Login:        user.Login,
		UserID:       user.ID,
	}, nil
}

// Validate validates the token.
func (s *Service) Validate(_ context.Context, tokenStr string) (*models.Token, error) {
	var err error
	if token, err := parseJWT(tokenStr, s.accessSecret); err == nil {
		if claims, ok := token.Claims.(*claims); ok && token.Valid {
			return &models.Token{
				Login:  claims.Login,
				UserID: claims.UserID,
			}, nil
		}
	}

	return nil, err
}

// parseJWT parses, validates, verifies the token with the specified cryptographic key
// for verifying the signature and returns the parsed token.
func parseJWT(tokenStr string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if s, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || s.Hash != crypto.SHA256 {
			return nil, ErrTokenBadSigningMethod
		}
		return secret, nil
	})
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, ErrInvalidToken
	}
	if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, ErrTokenExpired
	}

	return token, err
}

// Refresh refreshes the token.
func (s *Service) Refresh(ctx context.Context, tokenStr string) (*models.Token, error) {
	var err error
	if token, err := parseJWT(tokenStr, s.refreshSecret); err == nil {
		if claims, ok := token.Claims.(*claims); ok && token.Valid {
			user, err := s.store.User().FindByID(ctx, claims.UserID)
			if err != nil {
				return nil, err
			}

			return s.Create(ctx, user)
		}
	}

	return nil, err
}
