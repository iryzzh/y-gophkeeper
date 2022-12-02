package client

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

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
					Name:    "user",
					Aliases: []string{"u"},
					Usage:   "User",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
			},
		},
		{
			Name:   "auth",
			Usage:  "Log on to a remote server",
			Action: c.auth,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "remote",
					Aliases: []string{"r"},
					Usage:   "Remote server to authenticate with",
				},
				&cli.StringFlag{
					Name:    "user",
					Aliases: []string{"u"},
					Usage:   "User",
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
			Before: c.isInitialized,
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
					Name: "type",
					Usage: "Type of the entry." +
						" Should be one of: text, file or image",
					Value:   "text",
					Aliases: []string{"t"},
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:   "card",
					Usage:  "add a new bank card",
					Before: c.isInitialized,
					Action: c.entryNewCard,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "name",
							Aliases: []string{"n"},
							Usage:   "Name of the entry",
						},
						&cli.StringFlag{
							Name:    "type",
							Aliases: []string{"t"},
							Usage:   "Type of the card, e.g. 'Visa'",
						},
						&cli.StringFlag{
							Name:    "number",
							Aliases: []string{"cn"},
							Usage:   "Number of the card",
						},
						&cli.StringFlag{
							Name:    "exp",
							Aliases: []string{"e"},
							Usage:   "Expiration date MM/DD, e.g. '12/23'",
						},
						&cli.StringFlag{
							Name:  "cvv",
							Usage: "CVV",
						},
					},
				},
			},
		},
		{
			Name:   "view",
			Usage:  "View entries",
			Action: c.entryView,
			Before: c.isInitialized,
		},
		{
			Name:    "list",
			Aliases: []string{"ls", "l"},
			Usage:   "List entries",
			Action:  c.entryList,
			Before:  c.isInitialized,
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Usage:   "Delete entries",
			Action:  c.entryDelete,
			Before:  c.isInitialized,
		},
		{
			Name:    "version",
			Aliases: []string{"ver"},
			Usage:   "Show version",
			Action: func(context *cli.Context) error {
				fmt.Printf(
					"%s has version %s built from %s on %s",
					c.app.Name,
					c.cfg.Version.Version,
					c.cfg.Version.Commit,
					c.cfg.Version.BuildDate,
				)
				return nil
			},
		},
	}
}
