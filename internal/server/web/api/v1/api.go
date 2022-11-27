package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/iryzzh/gophkeeper/internal/services/item"
	"golang.org/x/net/context"

	"github.com/go-chi/chi/v5"
	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/services/token"
	"github.com/iryzzh/gophkeeper/internal/services/user"
	"github.com/pkg/errors"
)

type contextKey int
type empty struct{}

const (
	ctxUserID contextKey = iota
	ctxPageID
	ctxItem
)

// API is a http api service.
type API struct {
	tokenSvc *token.Service
	userSvc  *user.Service
	itemSvc  *item.Service
}

// NewAPI creates a new API.
func NewAPI(tokenSvc *token.Service, userSvc *user.Service, itemSvc *item.Service) *API {
	return &API{
		tokenSvc: tokenSvc,
		userSvc:  userSvc,
		itemSvc:  itemSvc,
	}
}

// Register registers the routes.
func (a *API) Register(r *chi.Mux) {
	r.Route("/api/v1/", func(r chi.Router) {
		r.Post("/signup", a.signup)
		r.Post("/login", a.login)
		r.Post("/token/refresh", a.tokenRefresh)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(a.Auth)
			r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			r.Route("/item", func(r chi.Router) {
				r.With(pagination).Get("/", a.itemGet)
				r.Get("/{id}", a.itemGet)
				r.With(itemCtx).Put("/", a.itemNew)
				r.With(itemCtx).Post("/{id}", a.itemSet)
				r.With(itemCtx).Delete("/{id}", a.itemSet)
			})
		})
	})
}

func itemCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var it *models.Item
		if err := json.NewDecoder(r.Body).Decode(&it); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), empty{}, nil)
		if usedID, ok := r.Context().Value(ctxUserID).(string); ok {
			it.UserID = usedID
			ctx = context.WithValue(r.Context(), ctxItem, it)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		intPage := 1
		if page != "" {
			var err error
			intPage, err = strconv.Atoi(page)
			if err != nil {
				http.Error(w, fmt.Sprintf("couldn't read page '%s': invalid syntax", page),
					http.StatusBadRequest)
				return
			}
		}
		ctx := context.WithValue(r.Context(), ctxPageID, intPage)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// signup matches the received login/password pair with the
// `models.User` struct, generates an encrypted password, creates
// the user in the database and, if successful, returns a new `models.Token`.
func (a *API) signup(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	create, code, err := a.userSvc.Create(r.Context(), u.Login, u.Password)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	t, err := a.tokenSvc.Create(r.Context(), create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, t, http.StatusCreated)
}

// login matches the received login/password pair with the
// `models.User` struct, validates it and, if successful, returns a
// new `models.Token`.
func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	login, err := a.userSvc.Login(r.Context(), u.Login, u.Password)
	if errors.Is(err, user.ErrInvalidUser) || errors.Is(err, user.ErrPasswordCannotBeEmpty) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, user.ErrLoginOrPasswordIsInvalid) || err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	t, createErr := a.tokenSvc.Create(r.Context(), login)
	if createErr != nil {
		http.Error(w, createErr.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, t, http.StatusOK)
}

// tokenRefresh matches the received token in `models.Token` format,
// validates it and, if successful, returns a new `models.Token`.
func (a *API) tokenRefresh(w http.ResponseWriter, r *http.Request) {
	var t models.Token
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newToken, err := a.tokenSvc.Refresh(r.Context(), t.RefreshToken)
	if errors.Is(err, token.ErrTokenExpired) || errors.Is(err, token.ErrInvalidToken) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("api sent: %v", newToken)

	WriteJSON(w, newToken, http.StatusCreated)
}

// itemResponse is used in the 'API.itemGet' handler
// responses, adding the 'meta.TotalPages' field to
// the 'models.Item' model. This field is used for
// pagination.
type itemResponse struct {
	Meta meta           `json:"meta,omitempty"`
	Data []*models.Item `json:"data"`
}

type meta struct {
	TotalPages int `json:"totalPages"`
}

// itemGet is a handler for incoming 'GET' requests to retrieve user
// `models.Item`.
//
// If a specific ID is specified as the path, the `models.Item` is
// returned if it exists.
//
// If a page number is specified as the query `?page=id`, the
// `[]models.Item` slice is returned.
func (a *API) itemGet(w http.ResponseWriter, r *http.Request) {
	if userID, ok := r.Context().Value(ctxUserID).(string); ok {
		if chi.URLParam(r, "id") != "" {
			byID, err := a.itemSvc.FindByID(r.Context(), userID, chi.URLParam(r, "id"))
			if err != nil {
				if errors.Is(err, item.ErrItemNotFound) {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			WriteJSON(w, itemResponse{Data: []*models.Item{byID}}, http.StatusOK)
			return
		}

		page := r.Context().Value(ctxPageID)
		items, pages, err := a.itemSvc.FindByUserID(r.Context(), userID, page)
		if err == nil {
			WriteJSON(w, itemResponse{Meta: meta{TotalPages: pages}, Data: items}, http.StatusOK)
			return
		}
		if errors.Is(err, item.ErrItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}
	http.Error(w, "", http.StatusBadRequest)
}

func (a *API) itemNew(w http.ResponseWriter, r *http.Request) {
	it, _ := r.Context().Value(ctxItem).(*models.Item)
	if err := a.itemSvc.Create(r.Context(), it); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, it, http.StatusCreated)
}

func (a *API) itemSet(w http.ResponseWriter, r *http.Request) {
	it, _ := r.Context().Value(ctxItem).(*models.Item)
	switch r.Method {
	case http.MethodPost:
		err := a.itemSvc.Update(r.Context(), it)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		WriteJSON(w, it, http.StatusOK)
		return
	case http.MethodDelete:
		err := a.itemSvc.Delete(r.Context(), it)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
