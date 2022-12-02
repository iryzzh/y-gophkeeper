package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/iryzzh/y-gophkeeper/internal/rand"
	"github.com/stretchr/testify/assert"

	"github.com/iryzzh/y-gophkeeper/internal/services/item"

	"github.com/iryzzh/y-gophkeeper/internal/store/sqlite"

	"github.com/iryzzh/y-gophkeeper/internal/store"

	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/iryzzh/y-gophkeeper/internal/services/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/utils"
	"github.com/stretchr/testify/require"
)

func testStore(t *testing.T) *sqlite.Store {
	t.Helper()
	cfg, err := utils.TestConfig(t)
	if err != nil {
		t.Fatal(err)
	}
	st, err := sqlite.NewStore(cfg.DB.DSN, "../../../../../migrations")
	if err != nil {
		t.Fatal(err)
	}

	return st
}

func newTestServer(t *testing.T, tokenSvc *token.Service, userSvc *user.Service, itemSvc *item.Service) (*httptest.Server, error) {
	t.Helper()

	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}

	h := chi.NewMux()
	apiV1 := NewAPI(tokenSvc, userSvc, itemSvc)
	apiV1.Register(h)

	ts := httptest.NewUnstartedServer(h)
	_ = ts.Listener.Close()
	ts.Listener = l

	ts.Start()

	return ts, nil
}

func testService(t *testing.T) (tokenSvc *token.Service, userSvc *user.Service, itemSvc *item.Service, st store.Store) {
	t.Helper()

	cfg, err := utils.TestConfig(t)
	if err != nil {
		t.Fatal(err)
	}

	st = testStore(t)

	tokenSvc = token.NewService(
		st,
		cfg.Security.AtExpiresIn,
		cfg.Security.RtExpiresIn,
		[]byte(cfg.Security.AccessSecret),
		[]byte(cfg.Security.RefreshSecret),
	)

	userSvc = user.NewService(
		st,
		cfg.Security.HashMemory,
		cfg.Security.HashIterations,
		cfg.Security.HashParallelism,
		cfg.Security.SaltLength,
		cfg.Security.KeyLength,
	)

	itemSvc = item.NewService(st)

	return tokenSvc, userSvc, itemSvc, st
}

func TestAPI_SignUp(t *testing.T) {
	tests := []struct {
		name string
		user *models.User
		want int
	}{
		{
			name: "ok",
			user: &models.User{
				Login:    "test",
				Password: "test",
			},
			want: http.StatusCreated,
		},
		{
			name: "conflict",
			user: &models.User{
				Login:    "test",
				Password: "test",
			},
			want: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			url := fmt.Sprintf("%v/api/v1/signup", ts.URL)

			if tt.want == http.StatusConflict {
				_ = st.User().Create(context.Background(), tt.user)
			}

			b, err := tt.user.Marshal()
			require.NoError(t, err)

			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetBody(b).
				Post(url)
			require.NoError(t, err)
			require.Equal(t, tt.want, resp.StatusCode())
		})
	}
}

func TestAPI_login(t *testing.T) {
	tests := []struct {
		name      string
		want      int
		runBefore func(svc *user.Service) error
		user      *models.User
	}{
		{
			name: "ok",
			user: &models.User{
				Login:    "test",
				Password: "test",
			},
			runBefore: func(svc *user.Service) error {
				err := svc.Create(context.Background(), &models.User{
					Login:    "test",
					Password: "test",
				})
				return err
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			url := fmt.Sprintf("%v/api/v1/login", ts.URL)

			require.NoError(t, tt.runBefore(uSvc))

			b, err := tt.user.Marshal()
			require.NoError(t, err)

			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetBody(b).
				Post(url)
			require.NoError(t, err)
			require.Equal(t, tt.want, resp.StatusCode())
		})
	}
}

func TestAPI_tokenRefresh(t *testing.T) {
	tests := []struct {
		name      string
		want      int
		runBefore func(uSvc *user.Service, tSvc *token.Service) (string, error)
	}{
		{
			name: "ok",
			runBefore: func(uSvc *user.Service, tSvc *token.Service) (string, error) {
				u := &models.User{Login: "test", Password: "test"}
				err := uSvc.Create(context.Background(), u)
				if err != nil {
					return "", err
				}
				tk, err := tSvc.Create(context.Background(), u)
				if err != nil {
					return "", err
				}

				return tk.RefreshToken, nil
			},
			want: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			url := fmt.Sprintf("%v/api/v1/token/refresh", ts.URL)

			rt, err := tt.runBefore(uSvc, tSvc)
			require.NoError(t, err)

			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetBody(models.Token{RefreshToken: rt}).
				Post(url)
			require.NoError(t, err)
			require.Equal(t, tt.want, resp.StatusCode())
		})
	}
}

func setupTestUserWithToken(t *testing.T, uSvc *user.Service, tSvc *token.Service) *models.Token {
	t.Helper()

	u := &models.User{Login: "test", Password: "test"}

	err := uSvc.Create(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := tSvc.Create(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}

	return tk
}

func TestAPI_itemGet(t *testing.T) {
	tests := []struct {
		name      string
		wantCode  int
		wantErr   error
		runBefore func(st store.Store, tSvc *token.Service, uSvc *user.Service) ([]*models.Item, string)
	}{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) ([]*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)
				testItem := utils.TestItem(t, withToken.UserID)
				err := st.Item().Create(context.Background(), testItem)
				if err != nil {
					t.Fatal(err)
				}
				return []*models.Item{testItem}, withToken.AccessToken
			},
		},
		{
			name:     "pagination",
			wantCode: http.StatusOK,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) ([]*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)
				var items []*models.Item
				for i := 0; i < 11; i++ {
					testItem := utils.TestItem(t, withToken.UserID)
					err := st.Item().Create(context.Background(), testItem)
					if err != nil {
						t.Fatal(err)
					}
					items = append(items, testItem)
				}

				return items, withToken.AccessToken
			},
		},
		{
			name:     "not found",
			wantCode: http.StatusNoContent,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) ([]*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)
				return []*models.Item{utils.TestItem(t, withToken.UserID)}, withToken.AccessToken
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			url := fmt.Sprintf("%v/api/v1/item", ts.URL)

			wantItems, accessToken := tt.runBefore(st, tSvc, uSvc)
			var itemsTotal []*models.Item

			client := resty.New()
			client.SetHeader("Accept", "application/json")
			client.SetAuthToken(accessToken)
			for i := 0; i < len(wantItems); i++ {
				if i%10 == 0 {
					client.SetQueryParams(map[string]string{
						"limit":  "10",
						"offset": fmt.Sprintf("%d", i),
					})
					resp, err := client.R().Get(url)
					require.NoError(t, err)
					require.Equal(t, tt.wantCode, resp.StatusCode())
					if resp.StatusCode() != http.StatusOK {
						return
					}
					got := &models.Items{}
					if err := json.Unmarshal(resp.Body(), &got); err != nil {
						t.Fatal(err)
					}
					itemsTotal = append(itemsTotal, got.Data...)
				}
			}
			require.Condition(t, func() bool {
				if tt.wantCode == http.StatusOK {
					return assert.Equal(t, wantItems, itemsTotal)
				}
				return true
			})
		})
	}
}

func TestAPI_itemGetWithID(t *testing.T) {
	tests := []struct {
		name      string
		wantCode  int
		runBefore func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string)
	}{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)
				testItem := utils.TestItem(t, withToken.UserID)
				err := st.Item().Create(context.Background(), testItem)
				if err != nil {
					t.Fatal(err)
				}

				return testItem, withToken.AccessToken
			},
		},
		{
			name:     "not found",
			wantCode: http.StatusNotFound,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)

				return &models.Item{ID: 999}, withToken.AccessToken
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			wantItem, accessToken := tt.runBefore(st, tSvc, uSvc)
			url := fmt.Sprintf("%v/api/v1/item/%v", ts.URL, wantItem.ID)
			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetAuthToken(accessToken).
				Get(url)
			require.NoError(t, err)
			require.Equal(t, tt.wantCode, resp.StatusCode())
			if tt.wantCode == http.StatusOK {
				var got models.Item
				if err := json.Unmarshal(resp.Body(), &got); err != nil {
					t.Fatal(err)
				}
				require.Equal(t, wantItem, &got)
			}
		})
	}
}

func TestAPI_itemNew(t *testing.T) {
	tests := []struct {
		name      string
		wantCode  int
		runBefore func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string)
	}{
		{
			name:     "created",
			wantCode: http.StatusCreated,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string) {
				return &models.Item{Meta: rand.String(10)}, setupTestUserWithToken(t, uSvc, tSvc).AccessToken
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			gotItem, accessToken := tt.runBefore(st, tSvc, uSvc)
			bytes, err := gotItem.Marshal()
			require.NoError(t, err)

			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetAuthToken(accessToken).
				SetBody(bytes).
				Put(fmt.Sprintf("%v/api/v1/item", ts.URL))
			require.NoError(t, err)
			require.Equal(t, tt.wantCode, resp.StatusCode())
			require.Condition(t, func() bool {
				if tt.wantCode == http.StatusCreated {
					ss := regexp.MustCompile(`"meta":"(.*?)"`).FindStringSubmatch(resp.String())
					if len(ss) < 1 {
						return false
					}
					return ss[1] == gotItem.Meta
				}
				return true
			})
		})
	}
}

func TestAPI_itemSetPost(t *testing.T) {
	tests := []struct {
		name      string
		wantCode  int
		runBefore func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string)
	}{
		{
			name:     "meta only",
			wantCode: http.StatusOK,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)
				testItem := &models.Item{
					UserID: withToken.UserID,
					Meta:   rand.String(10),
				}
				err := st.Item().Create(context.Background(), testItem)
				if err != nil {
					t.Fatal(err)
				}
				return &models.Item{ID: testItem.ID, Meta: rand.String(10)}, withToken.AccessToken
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			gotItem, accessToken := tt.runBefore(st, tSvc, uSvc)
			bytes, err := gotItem.Marshal()
			require.NoError(t, err)

			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetAuthToken(accessToken).
				SetBody(bytes).
				Post(fmt.Sprintf("%v/api/v1/item/%v", ts.URL, gotItem.ID))
			require.NoError(t, err)
			require.Equal(t, tt.wantCode, resp.StatusCode())
		})
	}
}

func TestAPI_itemSetDelete(t *testing.T) {
	tests := []struct {
		name      string
		wantCode  int
		runBefore func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string)
	}{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			runBefore: func(st store.Store, tSvc *token.Service, uSvc *user.Service) (*models.Item, string) {
				withToken := setupTestUserWithToken(t, uSvc, tSvc)
				testItem := &models.Item{
					UserID: withToken.UserID,
					Meta:   rand.String(10),
				}
				err := st.Item().Create(context.Background(), testItem)
				if err != nil {
					t.Fatal(err)
				}
				return &models.Item{ID: testItem.ID, Meta: rand.String(10)}, withToken.AccessToken
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tSvc, uSvc, iSvc, st := testService(t)
			ts, err := newTestServer(t, tSvc, uSvc, iSvc)
			require.NoError(t, err)
			defer func() {
				ts.Close()
				_ = st.Close()
			}()

			gotItem, accessToken := tt.runBefore(st, tSvc, uSvc)
			bytes, err := gotItem.Marshal()
			require.NoError(t, err)

			client := resty.New()
			resp, err := client.R().
				SetHeader("Accept", "application/json").
				SetAuthToken(accessToken).
				SetBody(bytes).
				Delete(fmt.Sprintf("%v/api/v1/item/%v", ts.URL, gotItem.ID))
			require.NoError(t, err)
			require.Equal(t, tt.wantCode, resp.StatusCode())
		})
	}
}
