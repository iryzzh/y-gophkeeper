package api_client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/iryzzh/y-gophkeeper/internal/config"
	"github.com/iryzzh/y-gophkeeper/internal/models"
)

const (
	apiPingEndpoint         = "/api/v1/ping"
	apiSignupEndpoint       = "/api/v1/signup"
	apiLoginEndpoint        = "/api/v1/login"
	apiRefreshTokenEndpoint = "/api/v1/token/refresh" //nolint:gosec
	apiItemEndpoint         = "/api/v1/item"
)

// ApiClient is a rest client.
type ApiClient struct {
	cfg   *config.API
	resty *resty.Client
}

// NewApiClient creates a new API client.
func NewApiClient(cfg *config.API, skipVerify bool) *ApiClient {
	ac := &ApiClient{
		cfg: cfg,
	}

	ac.resty = resty.New()
	ac.resty.SetBaseURL(cfg.Remote)
	ac.resty.SetAuthToken(cfg.AT)
	ac.resty.SetHeader("Accept", "application/json")
	ac.resty.SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipVerify, //nolint:gosec
		},
	})

	return ac
}

// Ping sends a ping request.
func (ac *ApiClient) Ping() error {
	get, err := ac.resty.R().Get(apiPingEndpoint)
	if err != nil {
		return err
	}

	if get.StatusCode() != http.StatusOK {
		return fmt.Errorf("remote ping failed")
	}

	return nil
}

// SetBaseURL sets the remote server address.
func (ac *ApiClient) SetBaseURL(url string) {
	ac.cfg.Remote = url
	ac.resty.SetBaseURL(url)
}

// Signup sends a signup request.
func (ac *ApiClient) Signup(user *models.User) error {
	body, err := user.Marshal()
	if err != nil {
		return err
	}

	resp, err := ac.resty.R().SetBody(body).Post(apiSignupEndpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("remote signup failed: %v", resp.String())
	}

	return ac.parseToken(resp.String())
}

func (ac *ApiClient) parseToken(token string) error {
	t := &models.Token{}
	if err := t.UnmarshalFromString(token); err != nil {
		return err
	}

	ac.cfg.AT = t.AccessToken
	ac.cfg.RT = t.RefreshToken

	ac.resty.SetAuthToken(t.AccessToken)

	return nil
}

func (ac *ApiClient) RefreshToken() error {
	resp, err := ac.resty.R().SetBody(models.Token{RefreshToken: ac.cfg.RT}).Post(apiRefreshTokenEndpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("remote refresh token failed: %v", resp.String())
	}

	return ac.parseToken(resp.String())
}

// Login sends a login request.
func (ac *ApiClient) Login(user *models.User) error {
	body, err := user.Marshal()
	if err != nil {
		return err
	}

	resp, err := ac.resty.R().SetBody(body).Post(apiLoginEndpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("remote login failed: %v", resp.String())
	}

	return ac.parseToken(resp.String())
}

type ActionKey string

const (
	ActionNew    ActionKey = "new"
	ActionUpdate ActionKey = "update"
	ActionDelete ActionKey = "delete"
)

func (ac *ApiClient) GetItems() ([]*models.Item, error) {
	var itemsTotal []*models.Item
	for i := 0; ; i++ {
		if i%10 == 0 {
			ac.resty.SetQueryParams(map[string]string{
				"limit":  "10",
				"offset": fmt.Sprintf("%d", i),
			})
			resp, err := ac.resty.R().Get(apiItemEndpoint)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode() == http.StatusNoContent {
				return itemsTotal, nil
			}
			if resp.StatusCode() != http.StatusOK {
				return nil, fmt.Errorf("remote get items failed: %v", resp.String())
			}
			got := &models.Items{}
			if err := json.Unmarshal(resp.Body(), &got); err != nil {
				return nil, err
			}
			itemsTotal = append(itemsTotal, got.Data...)
			if got.Meta.TotalItems == i {
				break
			}
		}
	}

	return itemsTotal, nil
}

func (ac *ApiClient) Item(item *models.Item, action ActionKey) error {
	body, err := item.Marshal()
	if err != nil {
		return err
	}

	switch action {
	case ActionNew:
		resp, err := ac.resty.R().SetBody(body).Put(apiItemEndpoint)
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusCreated {
			return fmt.Errorf("remote add item failed: %v", resp.String())
		}
	case ActionUpdate:
		resp, err := ac.resty.R().SetBody(body).Post(fmt.Sprintf("%v/%v", apiItemEndpoint, item.ID))
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("remote update item failed: %v", resp.String())
		}
	case ActionDelete:
		resp, err := ac.resty.R().SetBody(body).Delete(fmt.Sprintf("%v/%v", apiItemEndpoint, item.ID))
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("remote delete item failed: %v", resp.String())
		}
	}

	return nil
}
