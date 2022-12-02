package client

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/iryzzh/y-gophkeeper/internal/services/api_client"
	"github.com/iryzzh/y-gophkeeper/internal/services/item"
	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (c *Client) entryDelete(cCtx *cli.Context) error {
	name := cCtx.Args().First()

	userID, err := token.ParseUserIDFromToken(c.cfg.API.AT)
	if err != nil {
		return err
	}

	found, err := c.itemSvc.FindByMetaName(cCtx.Context, userID, name)
	if errors.Is(err, item.ErrItemNotFound) {
		return cli.Exit(fmt.Sprintf("entry '%v' not found", name), 1)
	}
	if err != nil {
		return err
	}

	if err = c.itemSvc.Delete(cCtx.Context, found); err != nil {
		color.Red("❌ item deletion failed: %v", err)
		return cli.Exit("", 1)
	}

	if err = c.clientSvc.RefreshToken(); err != nil {
		return err
	}

	if err = c.clientSvc.Item(found, api_client.ActionDelete); err != nil {
		return err
	}

	color.Green("✅ item was successfully deleted!")

	return nil
}
