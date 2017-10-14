// +build !appengine

package classic

import (
	"net/http"

	"golang.org/x/net/context"
)

func (c *Client) getClient(ctx context.Context) *http.Client {
	if c.HttpClient != nil {
		return c.HttpClient
	}
	return &http.Client{}
}
