package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClassic(t *testing.T) {
	a, err := NewClassic("demo", "2406471784", "demouser", "demopass")
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, a)
	assert.Equal(t, a.OrganizationID, "demo")
	assert.NotEqual(t, "", a.token)
}

func TestNewClassicErr(t *testing.T) {
	a, err := NewClassic("foo", "2406471784", "foouser", "foopass")
	assert.NotNil(t, a)
	assert.EqualError(t, err, "error: 401 Unauthorized")
}

// func TestNewClassicWithClientIDAndSecret(t *testing.T) {
// 	a := NewClassicWithClientIDAndSecret("demo", "foo", "bar")
// 	assert.Equal(t, a.OrganizationID, "demo")
// 	assert.Equal(t, a.ClientID, "foo")
// 	assert.Equal(t, a.ClientSecret, "bar")
// }
//
// func TestAPIMakeURL(t *testing.T) {
// 	a := &API{OrganizationID: "demo"}
// 	assert.Equal(t, "https://demo.memberclicks.net/foo/bar", a.makeURL("/foo/bar"))
// 	assert.Equal(t, "https://demo.memberclicks.net/foo/bar", a.makeURL("foo/bar"))
// 	assert.Equal(t, "https://demo.memberclicks.net/foo/bar", a.makeURL("https://demo.memberclicks.net/foo/bar"))
// }
//
// func TestAPIGetEndpoint(t *testing.T) {
// 	a := &API{OrganizationID: "demo"}
// 	assert.Equal(t, "https://demo.memberclicks.net", a.getEndpoint())
// }
//
// func TestRequestsSetAccessToken(t *testing.T) {
// 	a := &API{OrganizationID: "demo", AccessToken: &AccessToken{AccessToken: "abc"}}
// 	a.Post(context.Background(), "/api/v1/profile/me", url.Values{"foo": {"bar"}}, nil)
// 	if assert.NotNil(t, a.lastRequest) {
// 		assert.NotNil(t, a.lastResponse)
// 		assert.Equal(t, "application/x-www-form-urlencoded", a.lastRequest.Header.Get("Content-Type"))
// 		assert.Equal(t, "Bearer abc", a.lastRequest.Header.Get("Authorization"))
// 		assert.Equal(t, "POST", a.lastRequest.Method)
// 	}
// }
//
// func TestRequestSetAuthToken(t *testing.T) {
// 	a := &API{OrganizationID: "demo", authToken: "auth-token-foo-bar"}
// 	a.Get(context.Background(), "/api/v1/profile/me", nil)
// 	if assert.NotNil(t, a.lastRequest) {
// 		assert.NotNil(t, a.lastResponse)
// 		assert.Equal(t, "Bearer auth-token-foo-bar", a.lastRequest.Header.Get("Authorization"))
// 	}
// }
//
// func TestRequestSetBasicAuth(t *testing.T) {
// 	a := &API{OrganizationID: "demo", ClientID: "foo", ClientSecret: "bar"}
// 	a.Get(context.Background(), "/api/v1/profile/me", nil)
// 	if assert.NotNil(t, a.lastRequest) {
// 		u, p, _ := a.lastRequest.BasicAuth()
// 		assert.NotNil(t, a.lastResponse)
// 		assert.Equal(t, "foo", u)
// 		assert.Equal(t, "bar", p)
// 	}
// }
//
// func TestGetRequest(t *testing.T) {
// 	a := &API{OrganizationID: "demo", AccessToken: &AccessToken{AccessToken: "abc"}}
// 	a.Get(context.Background(), "/api/v1/profile/me?foo=bar", nil)
// 	if assert.NotNil(t, a.lastRequest) {
// 		assert.Equal(t, "GET", a.lastRequest.Method)
// 		assert.Equal(t, "bar", a.lastRequest.FormValue("foo"))
// 		assert.NotNil(t, a.lastResponse)
// 	}
// }
