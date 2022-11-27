package client

import (
	"context"
	"fmt"
	"os"

	"github.com/iryzzh/gophkeeper/internal/utils"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/iryzzh/gophkeeper/internal/tui"
	"github.com/urfave/cli/v2"
)

const (
	text  = "text"
	image = "image"
	file  = "file"
	card  = "card"
)

func (c *Client) entryNew(cCtx *cli.Context) error {
	name := cCtx.String("name")
	value := cCtx.String("value")
	valueType := cCtx.String("type")

	var err error
	if name == "" && value == "" {
		name, value, valueType, err = ask()
		if err != nil {
			return err
		}
	}

	if name == "" || value == "" {
		return fmt.Errorf("usage: add <name> <value>")
	}

	item := &models.Item{
		UserID: c.userUUID,
		Meta:   name,
		ItemData: &models.ItemData{
			Data: []byte(value),
		},
		DataType: valueType,
	}

	return c.s.Item().Create(context.Background(), item)
}

func ask() (name, value, valueType string, err error) {
	name, err = tui.AskString("Name:", nil)
	if err != nil {
		return "", "", "", err
	}
	valueType, err = tui.AskSelect("Entry type:", entryTypes(), entryTypes()[0])
	if err != nil {
		return "", "", "", err
	}

	switch valueType {
	case text:
		value, err = tui.AskString("Value:", nil)
	case file, image:
		value, err = tui.AskFile("File:")
		if err != nil {
			return "", "", "", err
		}
		bytes, err := os.ReadFile(value)
		if err != nil {
			return "", "", "", err
		}
		result := map[string]interface{}{
			value: bytes,
		}
		r, err := utils.InterfaceToBase64(&result)

		return name, string(r), valueType, err
	case card:
		var card models.Card
		card, err = tui.AskCard()
		if err != nil {
			return "", "", "", err
		}

		var data []byte
		data, err = card.EncodeToBase64()

		return name, string(data), valueType, err
	default:
		return "", "", "", fmt.Errorf("invalid data type: %v", valueType)
	}

	return name, value, valueType, err
}

func entryTypes() []string {
	return []string{text, image, file, card}
}
