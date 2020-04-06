package client

import (
	"fmt"
	"net/http"
	"net/url"
)

type Args map[string]string

type Client struct {
	Base *url.URL
}

func New(rawBaseURL string) (*Client, error) {
	client := new(Client)
	parsed, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}
	client.Base = parsed
	return client, nil
}

func (c *Client) get(path string, version int, endpoint string, params Args) (*http.Response, error) {
	basePath := fmt.Sprintf("%s/%v/v%v/%v", c.Base, path, version, endpoint)
	req, err := http.NewRequest("GET", basePath, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for key, val := range params {
		query.Add(key, val)
	}

	req.URL.RawQuery = query.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("(client) failed to make request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("(client) bad status on request: %v", res.StatusCode)
	}

	return res, nil
}
