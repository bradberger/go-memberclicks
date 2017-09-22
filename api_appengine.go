// +build appengine

package memberclicks

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// Client returns an HTTP client
func (a *API) getClient(ctx context.Context) *http.Client {
	if a.Client != nil {
		return a.Client
	}
	ctx, _ = context.WithTimeout(ctx, a.getTimeout())
	return urlfetch.Client(ctx)
}
