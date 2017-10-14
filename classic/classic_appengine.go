// +build appengine

package classic

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func (c *Client) getClient(ctx context.Context) *http.Client {
	if c.HttpClient != nil {
		return c.HttpClient
	}
	ctx, _ = context.WithTimeout(ctx, 15*time.Second)
	return urlfetch.Client(ctx)
}
