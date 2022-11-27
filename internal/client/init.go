package client

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/iryzzh/gophkeeper/internal/models"

	"github.com/urfave/cli/v2"
)

func (c *Client) init(cCtx *cli.Context) error {
	//TODO: return an error if the db has already been initialised
	remote := cCtx.String("remote")
	login := cCtx.String("login")
	password := cCtx.String("password")

	if remote == "" || login == "" || password == "" {
		return fmt.Errorf("usage: init <remote> <login> <password>")
	}

	//TODO: validate remote address
	resp, respErr := restyPost(fmt.Sprintf("%v/%v", remote, signupURL), "", models.User{
		Login:    login,
		Password: password,
	})
	if respErr != nil {
		return respErr
	}

	if resp.StatusCode() == http.StatusConflict {
		return fmt.Errorf("init: user already exists")
	}

	if resp.StatusCode() == http.StatusCreated {
		return c.initConfig(resp.Body(), remote, login)
	}

	return fmt.Errorf("unexpected response from remote: %v, code: %v", resp.String(), resp.StatusCode())
}

func (c *Client) initConfig(in []byte, remote, login string) error {
	var token models.Token
	if err := c.json.Unmarshal(in, &token); err != nil {
		return err
	}

	config := &LocalConfig{
		Remote:       remote,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Login:        login,
	}
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	return os.WriteFile(c.configFile, yamlData, 0o600) //nolint:gomnd
}
