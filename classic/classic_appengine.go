// +build appengine

package memberclicks

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

var (
	httpClient = urlfetch.Client(context.Background())
)

func UseContext(ctx context.Context) {
	httpClient = urlfetch.Client(ctx)
}
