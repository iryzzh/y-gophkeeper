package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/iryzzh/gophkeeper/internal/models"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func (c *Client) login(cCtx *cli.Context) error {
	remote := cCtx.String("remote")
	login := cCtx.String("login")
	password := cCtx.String("password")

	if remote == "" && login == "" && password == "" {
		var err error
		remote, login, password, err = c.askForCredentials()
		if err != nil {
			return err
		}
	}

	if remote == "" || login == "" || password == "" {
		return fmt.Errorf("usage: login <remote> <login> <password>")
	}

	resp, err := restyPost(fmt.Sprintf("%v/%v", remote, loginURL), "", models.User{Login: login, Password: password})
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusOK {
		return c.initConfig(resp.Body(), remote, login)
	}

	return fmt.Errorf("login: unexpected response from remote: %v, code: %v", resp.String(), resp.StatusCode())
}

func (c *Client) buildSuggestion(key string) func(in string) []string {
	_, err := os.Stat(c.configFile)
	if err != nil || os.IsNotExist(err) {
		return nil
	}

	file, err := os.ReadFile(c.configFile)
	if err != nil {
		return nil
	}

	parsed := make(map[string]string)
	if err = yaml.Unmarshal(file, &parsed); err != nil {
		return nil
	}
	for k, v := range parsed {
		if k == key {
			return func(in string) []string {
				if strings.HasPrefix(v, in) {
					return []string{v}
				}
				return nil
			}
		}
	}

	return nil
}
