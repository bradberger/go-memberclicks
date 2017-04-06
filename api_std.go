// +build !appengine

package memberclicks

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

var (
	httpClient *http.Client
)

// getClient returns an HTTP client
func (a *API) getClient(ctx context.Context) *http.Client {
	if a.Client != nil {
		return a.Client
	}
	if httpClient != nil {
		return httpClient
	}
	return &http.Client{Timeout: time.Second * 10}
}
