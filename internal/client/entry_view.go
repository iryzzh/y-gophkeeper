package client

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/iryzzh/gophkeeper/internal/clip"
	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/tui"
	"github.com/iryzzh/gophkeeper/internal/utils"
	"github.com/urfave/cli/v2"
)

func (c *Client) entryView(cCtx *cli.Context) error {
	name := cCtx.Args().First()

	items, err := c.getAllItems(context.Background())
	if err != nil {
		return err
	}

	if name == "" {
		return askOne(items)
	}

	for i, v := range items {
		if v.Meta == name {
			switch v.DataType {
			case file, image:
				return askToSaveTheFile(items[i].ItemData.Data)
			case card:
				var card models.Card
				if err := utils.FromBinaryBase64(items[i].ItemData.Data, &card); err != nil {
					return err
				}
				fmt.Printf(
					"type: %s, number: %s, expiry: %d/%d, ccv: %d\n",
					card.Type,
					card.Number,
					card.ExpiryMonth,
					card.ExpiryYear,
					card.CCV,
				)
				return nil
			}

			value := utils.FromBase64(items[i].ItemData.Data)

			return maybeCopyToClipboard(value)
		}
	}

	return fmt.Errorf("entry not found")
}

func askOne(itemsTotal []*models.Item) error {
	var vs []tui.ViewStruct //nolint:prealloc
	for i, v := range itemsTotal {
		vs = append(vs, tui.ViewStruct{
			Title: v.Meta,
			ID:    i,
		})
	}
	idx, err := tui.View(vs)
	if err != nil {
		return err
	}

	v := utils.FromBase64(itemsTotal[idx].ItemData.Data)

	return maybeCopyToClipboard(v)
}

func askToSaveTheFile(src []byte) error {
	var m map[string]interface{}
	if err := utils.FromBinaryBase64(src, &m); err != nil {
		return err
	}
	for k, v := range m {
		s, _ := v.(string)

		split := strings.Split(k, string(os.PathSeparator))
		name := split[len(split)-1]
		fmt.Printf("Original file: %v\n", name)

		path, err := tui.AskFile("Save file to:", false)
		if err != nil {
			return err
		}

		return os.WriteFile(path, utils.FromBase64([]byte(s)), 0o600) //nolint:gomnd
	}

	return nil
}

func maybeCopyToClipboard(value []byte) error {
	fmt.Printf("value: %v\n", string(value))

	_ = clip.Write(value, 0)

	return nil
}
