// +build !appengine

package memberclicks

import (
	"net/http"

	"golang.org/x/net/context"
)

func (a *API) getClient(ctx context.Context) *http.Client {
	return http.DefaultClient
}
