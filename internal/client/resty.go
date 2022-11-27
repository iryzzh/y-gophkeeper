package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

func newClient(token string) *resty.Client {
	client := resty.New()
	client.SetTransport(customTransport())
	client.SetHeader("Accept", "application/json")

	if len(token) > 0 {
		client.SetAuthToken(token)
	}

	return client
}

func restyGet(url, token string) (*resty.Response, error) {
	client := newClient(token)

	return client.R().Get(url)
}

func restyPost(url, token string, payload interface{}) (*resty.Response, error) {
	client := newClient(token)

	return client.R().SetBody(payload).Post(url)
}

func restyPut(url, token string, payload interface{}) (*resty.Response, error) {
	client := newClient(token)

	return client.R().SetBody(payload).Put(url)
}

func restyDelete(url, token string, payload interface{}) (*resty.Response, error) {
	client := newClient(token)

	return client.R().SetBody(payload).Delete(url)
}

//nolint:gomnd
func customTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		},
	}
}
