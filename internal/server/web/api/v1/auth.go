package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/iryzzh/y-gophkeeper/internal/models"
	"golang.org/x/net/context"
)

// Auth is an authentication middleware.
func (a *API) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := a.verifyRequest(r, tokenFromHeader)
		if err != nil || token == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, ctxUserID, token.UserID)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) verifyRequest(r *http.Request, findTokenFns ...func(r *http.Request) string) (*models.Token, error) {
	var tokenString string

	for _, fn := range findTokenFns {
		tokenString = fn(r)
		if tokenString != "" {
			break
		}
	}

	if tokenString == "" {
		return nil, fmt.Errorf("token not found")
	}

	return a.tokenSvc.Validate(r.Context(), tokenString)
}

// // tokenFromCookie tries to retrieve the token string from a cookie named
// // "jwt".
// func tokenFromCookie(r *http.Request) string {
//	cookie, err := r.Cookie("jwt")
//	if err != nil {
//		return ""
//	}
//	return cookie.Value
// }

// tokenFromHeader tries to retrieve the token string from the
// "Authorization" request header: "Authorization: BEARER T".
func tokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
