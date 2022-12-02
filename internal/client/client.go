package client

import (
	"context"
	"os"

	"github.com/iryzzh/y-gophkeeper/internal/config"
	"github.com/iryzzh/y-gophkeeper/internal/services/api_client"
	"github.com/iryzzh/y-gophkeeper/internal/services/item"
	"github.com/iryzzh/y-gophkeeper/internal/services/user"
	"github.com/iryzzh/y-gophkeeper/internal/store"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
)

type Client struct {
	json      jsoniter.API
	store     store.Store
	app       *cli.App
	cfg       *config.ClientCfg
	userSvc   *user.Service
	itemSvc   *item.Service
	clientSvc *api_client.ApiClient
}

func NewClient(cfg *config.ClientCfg, s store.Store) *Client {
	c := &Client{
		cfg:   cfg,
		store: s,
		json:  jsoniter.ConfigCompatibleWithStandardLibrary,
	}

	c.userSvc = user.NewService(
		s,
		cfg.Security.HashMemory,
		cfg.Security.HashIterations,
		cfg.Security.HashParallelism,
		cfg.Security.SaltLength,
		cfg.Security.KeyLength,
	)

	c.itemSvc = item.NewService(s)

	c.clientSvc = api_client.NewApiClient(&cfg.API, cfg.SkipVerify)

	commands := c.getCommands()
	c.app = &cli.App{
		Commands: commands,
	}
	c.app.Name = "gophkeeper-cli"
	c.app.Version = cfg.Version.Version

	return c
}

func (c *Client) Run(ctx context.Context) error {
	return c.app.RunContext(ctx, os.Args)
}

func (c *Client) updateToken() error {
	if err := c.clientSvc.RefreshToken(); err != nil {
		return err
	}

	return c.cfg.SaveConfig()
}
