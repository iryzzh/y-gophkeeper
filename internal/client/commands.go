package client

import (
	"fmt"

	"github.com/iryzzh/gophkeeper/internal/tui"

	"github.com/urfave/cli/v2"
)

func (c *Client) isInitialised(_ *cli.Context) error {
	if c.userUUID != "" {
		return nil
	}

	fmt.Print(logo)
	fmt.Sprintln("🌟 Welcome to Gophkeeper!")
	fmt.Println("No existing configuration found.")
	fmt.Println("☝ Please run 'gophkeeper init'")

	return fmt.Errorf("not initialised")
}

func (c *Client) askForCredentials() (remote, login, password string, err error) {
	remote, err = tui.AskString("remote:", c.buildSuggestion("remote"))
	if err != nil {
		return "", "", "", err
	}
	login, err = tui.AskString("login:", c.buildSuggestion("login"))
	if err != nil {
		return "", "", "", err
	}
	password, err = tui.AskPassword()
	if err != nil {
		return "", "", "", err
	}

	return remote, login, password, nil
}

func (c *Client) getCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize a new gophkeeper",
			Action: c.init,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "remote",
					Aliases: []string{"r"},
					Usage:   "Remote server for data synchronisation",
				},
				&cli.StringFlag{
					Name:    "login",
					Aliases: []string{"l"},
					Usage:   "Login",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
			},
		},
		{
			Name:   "login",
			Usage:  "Authenticate to a remote",
			Action: c.login,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "remote",
					Aliases: []string{"r"},
					Usage:   "Remote server to authenticate with",
				},
				&cli.StringFlag{
					Name:    "login",
					Aliases: []string{"l"},
					Usage:   "Login",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
			},
		},
		{
			Name:   "add",
			Usage:  "Add an entry",
			Action: c.entryNew,
			Before: c.isInitialised,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Usage:   "Name of the entry",
				},
				&cli.StringFlag{
					Name:    "value",
					Aliases: []string{"v"},
					Usage:   "Value of the entry",
				},
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Usage: "Type of the entry." +
						" Should be one of: text, file, image or card (bank card)",
					Value: "text",
				},
			},
		},
		{
			Name:   "view",
			Usage:  "View entries",
			Action: c.entryView,
			Before: c.isInitialised,
		},
		{
			Name:    "list",
			Aliases: []string{"ls", "l"},
			Usage:   "List entries",
			Action:  c.entryList,
			Before:  c.isInitialised,
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Usage:   "Delete entries",
			Action:  c.entryDelete,
			Before:  c.isInitialised,
		},
		{
			Name:   "push",
			Usage:  "Push local items to a remote server",
			Action: c.push,
			Before: c.isInitialised,
		},
		{
			Name:   "pull",
			Usage:  "Retrieve records from a remote server",
			Action: c.pull,
			Before: c.isInitialised,
		},
		{
			Name:    "version",
			Aliases: []string{"ver"},
			Usage:   "Show version",
			Action: func(context *cli.Context) error {
				fmt.Println("show version is fired")
				return nil
			},
		},
	}
}
