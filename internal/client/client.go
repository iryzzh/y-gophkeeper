package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/store"

	"github.com/golang-jwt/jwt/v4"
	"github.com/iryzzh/gophkeeper/internal/store/sqlite"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	signupURL       = "api/v1/signup"
	loginURL        = "api/v1/login"
	itemURL         = "api/v1/item"
	tokenRefreshURL = "api/v1/token/refresh"
	pingURL         = "api/v1/ping"
)

type Client struct {
	app        *cli.App
	configFile string
	json       jsoniter.API
	s          *sqlite.Store
	userUUID   string
	cfg        LocalConfig
}

type LocalConfig struct {
	Remote       string `yaml:"remote"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
	Login        string `yaml:"login"`
}

func NewClient(_ context.Context, s *sqlite.Store) *Client {
	c := &Client{
		s:          s,
		configFile: "config.yml",
		json:       jsoniter.ConfigCompatibleWithStandardLibrary,
	}
	c.app = c.setupApp()

	return c
}

func (c *Client) Run() error {
	if err := c.ping(); err != nil {
		return err
	}

	return c.app.Run(os.Args)
}

func (c *Client) setupApp() *cli.App {
	c.readConfig()

	commands := c.getCommands()
	app := &cli.App{
		Commands: commands,
	}

	return app
}

func (c *Client) readConfig() {
	file, err := os.ReadFile(c.configFile)
	if err != nil {
		return
	}

	var config LocalConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return
	}

	c.cfg = config

	c.readUserUUID(config.AccessToken)
}

func (c *Client) readUserUUID(tokenStr string) {
	token, _ := jwt.Parse(tokenStr, nil)
	if token == nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for k, v := range claims {
			if k == "user_id" {
				if uuid, ok := v.(string); ok {
					c.userUUID = uuid
				}

				return
			}
		}
	}
}

func (c *Client) getAllItems(ctx context.Context) ([]*models.Item, error) {
	var itemsTotal []*models.Item
	for i := 1; ; i++ {
		items, page, err := c.s.Item().FindByUserID(ctx, c.userUUID, i)
		if err != nil {
			if errors.Is(err, store.ErrItemNotFound) {
				return nil, fmt.Errorf("no items")
			}
			return nil, err
		}
		itemsTotal = append(itemsTotal, items...)
		if i == page {
			break
		}
	}

	sort.Slice(itemsTotal, func(i, j int) bool {
		return itemsTotal[i].Meta < itemsTotal[j].Meta
	})
	return itemsTotal, nil
}

func (c *Client) ping() error {
	client := newClient(c.cfg.AccessToken)
	resp, err := client.R().Get(fmt.Sprintf("%v/%v", c.cfg.Remote, pingURL))
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return c.refreshToken()
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected response from remote: %v, code: %v", resp.String(), resp.StatusCode())
	}

	return nil
}

func (c *Client) refreshToken() error {
	client := newClient("")

	resp, err := client.R().SetBody(models.Token{RefreshToken: c.cfg.RefreshToken}).
		Post(fmt.Sprintf("%v/%v", c.cfg.Remote, tokenRefreshURL))
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d, msg: %v", resp.StatusCode(), resp.String())
	}

	fmt.Printf("send body: %v\n", resp.String())

	return c.initConfig(resp.Body(), c.cfg.Remote, c.cfg.Login)
}
