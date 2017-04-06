// +build !appengine

package memberclicks

import "net/http"

var (
	httpClient = http.DefaultClient
)
