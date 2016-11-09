// +build appengine

package memberclicks

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func (a *API) getClient(ctx context.Context) *http.Client {
	return urlfetch.Client(ctx)
}
