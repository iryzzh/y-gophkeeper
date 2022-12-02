package client

import (
	"strings"

	"github.com/fatih/color"
	"github.com/iryzzh/y-gophkeeper/internal/models"
	"github.com/iryzzh/y-gophkeeper/internal/services/api_client"
	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/iryzzh/y-gophkeeper/internal/tui"
	"github.com/urfave/cli/v2"
)

func (c *Client) entryNewCard(cCtx *cli.Context) (err error) {
	var meta string
	card := &models.Card{}
	if cCtx.NumFlags() != 0 {
		card.Type = cCtx.String("type")
		card.Number = cCtx.String("number")
		card.CVV = cCtx.String("cvv")

		split := strings.Split(cCtx.String("exp"), `/`)
		if len(split) > 1 {
			card.Month = split[0]
			card.Year = split[1]
		}

		meta = cCtx.String("name")
	} else {
		meta, err = tui.AskCard(card)
		if err != nil {
			return err
		}
	}

	if card == nil || meta == "" {
		return cli.Exit("usage: add card --name <name> --type <type> --number <number> --exp <expiration> --cvv <cvv>", 1)
	}

	if err = card.Validate(); err != nil {
		color.Red("❌ %v", err)
		return cli.Exit("", 1)
	}

	var data []byte
	if data, err = card.Encode(); err != nil {
		return err
	}

	var userID string
	userID, err = token.ParseUserIDFromToken(c.cfg.API.AT)
	if err != nil {
		return err
	}

	item := &models.Item{
		UserID:   userID,
		Meta:     meta,
		DataType: models.EntryTypeCard,
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

	color.Green("✅ card was successfully created!")

	return nil
}
