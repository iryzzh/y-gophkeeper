package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iryzzh/gophkeeper/internal/store"
	"github.com/pkg/errors"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/urfave/cli/v2"
)

func (c *Client) push(_ *cli.Context) error {
	items, err := c.getAllItems(context.Background())
	if err != nil {
		return err
	}

	if err := c.deleteRemoteAll(); err != nil {
		return err
	}

	if err := c.pushAllItems(items); err != nil {
		return err
	}

	fmt.Printf("done.\n")
	return nil
}

func (c *Client) pull(_ *cli.Context) error {
	items, err := c.getRemoteItems()
	if err != nil {
		return err
	}

	for _, v := range items {
		err := c.s.Item().Create(context.Background(), v)
		if errors.Is(err, store.ErrItemExists) {
			err = c.s.Item().Update(context.Background(), v)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("done.\n")
	return nil
}

func (c *Client) pushAllItems(items []*models.Item) error {
	for _, v := range items {
		url := fmt.Sprintf("%v/%v", c.cfg.Remote, itemURL)
		payload, err := json.Marshal(v)
		if err != nil {
			return err
		}

		resp, err := restyPut(url, c.cfg.AccessToken, payload)
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusCreated {
			return fmt.Errorf("unexpected response from remote: %v, code: %v", resp.String(), resp.StatusCode())
		}
	}

	return nil
}

func (c *Client) deleteRemoteAll() error {
	items, err := c.getRemoteItems()
	if err != nil {
		return err
	}

	for _, v := range items {
		url := fmt.Sprintf("%v/%v/%v", c.cfg.Remote, itemURL, v.ID)
		payload, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = restyDelete(url, c.cfg.AccessToken, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) getRemoteItems() ([]*models.Item, error) {
	var items []*models.Item
	for i := 1; ; i++ {
		rURL := fmt.Sprintf("%v/%v/?page=%v", c.cfg.Remote, itemURL, i)
		resp, err := restyGet(rURL, c.cfg.AccessToken)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode() == http.StatusNotFound {
			return items, nil
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("bad status code: %v", resp.StatusCode())
		}

		var got models.ItemResponse
		if err := json.Unmarshal(resp.Body(), &got); err != nil {
			return nil, err
		}

		items = append(items, got.Data...)

		if i == got.Meta.TotalPages {
			return items, nil
		}
	}
}
