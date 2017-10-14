package classic

import (
	"net/url"
	"os"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

var (
	orgID    = os.Getenv("MEMBERCLICKS_ORG_ID")
	apiKey   = os.Getenv("MEMBERCLICKS_API_KEY")
	username = os.Getenv("MEMBERCLICKS_USERNAME")
	password = os.Getenv("MEMBERCLICKS_PASSWORD")

	networkTests = false
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	a, err := New(ctx, orgID, apiKey, username, password)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, a)
	assert.Equal(t, a.OrganizationID, orgID)
	assert.NotEqual(t, "", a.token)
}

func TestNewClassicErr(t *testing.T) {
	ctx := context.Background()
	a, err := New(ctx, "foobar", apiKey, username, password)
	assert.NotNil(t, a)
	assert.EqualError(t, err, "error: 401 Unauthorized")
}

func TestGetUsers(t *testing.T) {
	ctx := context.Background()
	a, err := New(ctx, orgID, apiKey, username, password)
	if !assert.NoError(t, err) {
		return
	}
	params := url.Values{"pageNum": {"1"}, "pageSize": {"10"}}
	list, err := a.Users(ctx, params)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.True(t, len(list.Users) == 10) {
		return
	}
	u, err := a.User(ctx, list.Users[0].UserID)
	assert.NoError(t, err)
	if !assert.NotNil(t, u) {
		return
	}
	assert.Equal(t, u.UserID, list.Users[0].UserID)
}

func TestAPIMakeURL(t *testing.T) {
	a := &Client{OrganizationID: "demo"}
	assert.Equal(t, "https://demo.memberclicks.net/foo/bar", a.makeURL("/foo/bar"))
	assert.Equal(t, "https://demo.memberclicks.net/foo/bar", a.makeURL("foo/bar"))
	assert.Equal(t, "https://demo.memberclicks.net/foo/bar", a.makeURL("https://demo.memberclicks.net/foo/bar"))
}

func TestAPIGetEndpoint(t *testing.T) {
	a := &Client{OrganizationID: "demo"}
	assert.Equal(t, "https://demo.memberclicks.net", a.getEndpoint())
}
