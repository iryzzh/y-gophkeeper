package client

import (
	"fmt"
	"strings"

	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/urfave/cli/v2"
)

func (c *Client) entryList(cCtx *cli.Context) error {
	userID, err := token.ParseUserIDFromToken(c.cfg.API.AT)
	if err != nil {
		return err
	}

	items, err := c.itemSvc.FindByUserID(cCtx.Context, userID, "1000", "0")
	if err != nil {
		return err
	}

	var prev string
	for _, v := range items.Data {
		s := strings.Split(v.Meta, `/`)
		if s[0] != prev {
			fmt.Printf("%s\n", s[0])
		}
		for i := 1; i < len(s); i++ {
			printListing(s[i], i-1)
		}
		prev = s[0]
	}

	return nil
}

func printListing(entry string, depth int) {
	indent := strings.Repeat("│   ", depth)
	fmt.Printf("%s└── %s\n", indent, entry)
}
