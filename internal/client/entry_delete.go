package client

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *Client) entryDelete(cCtx *cli.Context) error {
	name := cCtx.Args().First()
	// if name == "" {
	//	//
	// }

	items, err := c.getAllItems(context.Background())
	if err != nil {
		return err
	}

	for _, v := range items {
		if v.Meta == name {
			if err := c.s.Item().Delete(context.Background(), v); err != nil {
				return err
			}

			fmt.Printf("item '%v' is deleted\n", name)
			return nil
		}
	}

	return nil
}
