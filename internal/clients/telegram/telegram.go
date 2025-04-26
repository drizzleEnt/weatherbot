package telegram

import "net/http"

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host string, token string) *Client{
	return &Client{
		host:     host,
		basePath: basePath(token),
		client:   http.Client{},
	}
}

func basePath(token string) string{
	return "bot" + token
}