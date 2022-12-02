package client

import (
	"fmt"
	"os"

	"github.com/iryzzh/y-gophkeeper/internal/config"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/services/item"
	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/iryzzh/y-gophkeeper/internal/tui"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (c *Client) entryView(cCtx *cli.Context) error {
	name := cCtx.Args().First()

	userID, err := token.ParseUserIDFromToken(c.cfg.API.AT)
	if err != nil {
		return err
	}

	var foundItem *models.Item
	foundItem, err = c.itemSvc.FindByMetaName(cCtx.Context, userID, name)
	if errors.Is(err, item.ErrItemNotFound) {
		return cli.Exit(fmt.Sprintf("entry '%v' not found", name), 1)
	}
	if err != nil {
		return err
	}

	var data []byte
	data, err = foundItem.ItemData.DecodeDataToBytes()
	if err != nil {
		return err
	}

	switch foundItem.DataType {
	case models.EntryTypeCard:
		card := &models.Card{}
		if err = card.Decode(foundItem.ItemData.Data); err != nil {
			return err
		}
		fmt.Printf(
			"type: %s, number: %s, expiry: %s/%s, cvv: %s\n",
			card.Type,
			card.Number,
			card.Month,
			card.Year,
			card.CVV,
		)

		return nil
	case models.EntryTypeImage, models.EntryTypeFile:
		var path string
		path, err = tui.AskFile("save the file to:", false)
		if err != nil {
			return err
		}

		return os.WriteFile(path, data, config.FilePermission)
	default:
		fmt.Printf("value: %v\n", data)
	}

	return nil
}
