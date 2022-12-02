package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/tui"
	"github.com/urfave/cli/v2"
)

// isInitialized checks for existing configuration.
func (c *Client) isInitialized(cCtx *cli.Context) error {
	printMsg := func() {
		executable, _ := os.Executable()
		base := filepath.Base(executable)
		var name string
		if base == "" {
			name = cCtx.App.Name
		} else {
			name = base
		}
		fmt.Printf("%v\n", logo)
		fmt.Printf("Existing configuration not found.\n")
		fmt.Printf("☝ Please run '%s init'\n", name)
	}

	initialized, err := c.store.IsUsersExist()
	if err != nil {
		return err
	}

	if !initialized {
		printMsg()
		return cli.Exit("", 1)
	}

	return nil
}

// init initializes the configuration.
func (c *Client) init(cCtx *cli.Context) error {
	initModel := &models.Init{
		Remote: cCtx.String("remote"),
		User: &models.User{
			Login:    cCtx.String("user"),
			Password: cCtx.String("password")},
	}

	fmt.Printf("%v\n", logo)
	fmt.Printf("Initializing a new %v\n", cCtx.App.Name)

	if !cCtx.IsSet("remote") && !cCtx.IsSet("user") && !cCtx.IsSet("password") {
		err := tui.AskInit(initModel)
		if err != nil {
			return err
		}
	}

	if initModel.Remote == "" || initModel.User.Login == "" || initModel.User.Password == "" {
		return cli.Exit("usage: init <remote> <auth> <password>", 1)
	}

	if err := c.initStore(cCtx.Context, initModel); err != nil {
		return err
	}

	color.Green("✅ initialization completed successfully!")
	return nil
}

func (c *Client) initStore(ctx context.Context, initModel *models.Init) error {
	usersExist, err := c.store.IsUsersExist()
	if err != nil {
		return err
	}

	if usersExist {
		color.Red("❌ The store has already been initialized")
		var confirm bool
		if err = tui.AskConfirm("continue?", &confirm); err != nil {
			return err
		}
		if !confirm {
			color.Yellow("canceled")
			return cli.Exit("", 1)
		}
	}

	c.clientSvc.SetBaseURL(initModel.Remote)

	if err = c.clientSvc.Signup(initModel.User); err != nil {
		return err
	}

	if err = c.userSvc.Create(ctx, initModel.User); err != nil {
		return err
	}

	return c.cfg.SaveConfig()
}
