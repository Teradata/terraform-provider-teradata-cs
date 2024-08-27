package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const HostURL string = "https://api.clearscape.teradata.com/"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(host, token string) (*Client, error) {
	c := Client{
		HostURL:    HostURL,
		HTTPClient: &http.Client{Timeout: 1000 * time.Second},
		Token:      token,
	}

	if host != "" {
		c.HostURL = host
	}

	c.Token = token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Token

	req.Header.Set("Authorization", "Bearer "+token)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}
	return body, err
}
