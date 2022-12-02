package client

import (
	"github.com/fatih/color"
	"github.com/iryzzh/y-gophkeeper/internal/file"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/services/api_client"
	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/iryzzh/y-gophkeeper/internal/tui"
	"github.com/urfave/cli/v2"
)

func (c *Client) entryNew(cCtx *cli.Context) error {
	entry := &models.Entry{
		Name:      cCtx.String("name"),
		Value:     cCtx.String("value"),
		EntryType: cCtx.String("type"),
	}
	color.Cyan("📝 creating a new entry")

	if !cCtx.IsSet("name") && !cCtx.IsSet("value") && !cCtx.IsSet("type") {
		if err := tui.AskEntry(entry); err != nil {
			return err
		}
	}

	userID, err := token.ParseUserIDFromToken(c.cfg.API.AT)
	if err != nil {
		return err
	}

	var data []byte
	if entry.EntryType == models.EntryTypeImage || entry.EntryType == models.EntryTypeFile {
		data, err = file.Encode(entry.Value)
		if err != nil {
			color.Red("❌ %v", err)
			return cli.Exit("", 1)
		}
	} else {
		data = entry.EncodeBytes()
	}

	item := &models.Item{
		UserID:   userID,
		Meta:     entry.Name,
		DataType: entry.EntryType,
		ItemData: &models.ItemData{
			Data: data,
		},
	}

	err = c.itemSvc.Create(cCtx.Context, item)
	if err != nil {
		color.Red("❌ %v", err)
		return cli.Exit("", 1)
	}

	if err = c.clientSvc.RefreshToken(); err != nil {
		return err
	}

	if err = c.clientSvc.Item(item, api_client.ActionNew); err != nil {
		return err
	}

	color.Green("✅ item was successfully created!")

	return nil
}
