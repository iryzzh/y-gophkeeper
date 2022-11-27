package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func (c *Client) entryList(_ *cli.Context) error {
	items, err := c.getAllItems(context.Background())
	if err != nil {
		return err
	}

	var prev string
	for _, v := range items {
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
