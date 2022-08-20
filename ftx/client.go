package ftx

import (
	"net/http"
	"net/url"
)

type Client struct {
	Client     *http.Client
	Api        string
	Secret     []byte
	Subaccount string

	markets map[string]Market
}

func New(api string, secret string, subaccount string) *Client {
	return &Client{Client: &http.Client{}, Api: api, Secret: []byte(secret), Subaccount: url.PathEscape(subaccount),
		markets: make(map[string]Market)}
}
