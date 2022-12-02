package client

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/tui"
	"github.com/urfave/cli/v2"
)

// auth authorizes on the remote server.
func (c *Client) auth(cCtx *cli.Context) error {
	initModel := &models.Init{
		Remote: cCtx.String("remote"),
		User: &models.User{
			Login:    cCtx.String("user"),
			Password: cCtx.String("password")},
	}

	fmt.Printf("%v\n", logo)

	if !cCtx.IsSet("remote") && !cCtx.IsSet("user") && !cCtx.IsSet("password") {
		err := tui.AskInit(initModel)
		if err != nil {
			return err
		}
	}

	if initModel.Remote == "" || initModel.User.Login == "" || initModel.User.Password == "" {
		return cli.Exit("usage: auth <remote> <auth> <password>", 1)
	}

	fmt.Printf("🗣️ logging on to %v...\n", initModel.Remote)

	c.clientSvc.SetBaseURL(initModel.Remote)
	if err := c.clientSvc.Login(initModel.User); err != nil {
		return err
	}

	_ = c.userSvc.Create(cCtx.Context, initModel.User)

	if err := c.cfg.SaveConfig(); err != nil {
		return err
	}

	color.Green("✅ success!")

	return nil
}
