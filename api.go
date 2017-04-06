package memberclicks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

// Scopes
const (
	ScopeRead = "read"
)

// API is an MemberClicks API client
type API struct {
	orgID, clientID, clientSecret, accessToken string

	Client *http.Client
}

// GetAuthCodeURL returns a auth code URL for redirecting the client to authorize with MemberClicks
func (a *API) GetAuthCodeURL(scope, state, redirectURL string) string {
	return fmt.Sprintf(
		"https://%s.memberclicks.net/oauth/v1/authorize?response_type=code&client_id=%s&scope=%s&state=%s&redirect_uri=%s",
		a.orgID,
		a.clientID,
		scope,
		state,
		redirectURL,
	)
}

// AuthCodeRedirect does an http.Redirect to the authcodeurl
func (a *API) AuthCodeRedirect(w http.ResponseWriter, r *http.Request, scope, state, redirectURL string) {
	http.Redirect(w, r, a.GetAuthCodeURL(scope, state, redirectURL), http.StatusTemporaryRedirect)
}

// GetToken trades an auth code for an access token
func (a *API) GetToken(ctx context.Context, authCode, scope, state, redirectURL string) (*Token, error) {

	var t Token
	form := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {authCode},
		"scope":        {scope},
		"state":        {state},
		"redirect_uri": {redirectURL},
	}
	if err := a.Post(ctx, fmt.Sprintf("https://%s.memberclicks.net/oauth/v1/token", a.orgID), form, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// SetAccessToken sets the internal access token to use on requests if no other Authorization header is set.
func (a *API) SetAccessToken(accessToken string) *API {
	a.accessToken = accessToken
	return a
}

// Auth initializes a default ClientCredentials requests and stores the resulting access token if successful
func (a *API) Auth(ctx context.Context) error {
	t, err := a.ClientCredentials(ctx, "read")
	if err != nil {
		return err
	}
	a.accessToken = t.AccessToken
	return nil
}

// ClientCredentials returns a client_credentials token
func (a *API) ClientCredentials(ctx context.Context, scope string) (*Token, error) {
	var t Token
	form := url.Values{"grant_type": {"client_credentials"}, "scope": {scope}}
	if err := a.Post(ctx, "/oauth/v1/token", form, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// RefreshToken gets a new token from a refresh token
func (a *API) RefreshToken(ctx context.Context, scope string, refreshToken string) (*Token, error) {
	var t Token
	form := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {refreshToken}}
	if err := a.Post(ctx, "/oauth/v1/token", form, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Me returns the profile associated with the accessToken
func (a *API) Me(ctx context.Context, accessToken string) (*Profile, error) {

	var p Profile
	req, err := http.NewRequest("GET", a.makeURL("/api/v1/profile/me"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err := a.Do(ctx, req, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Post sends a POST request to the urlStr and marshals the response into result
func (a *API) Post(ctx context.Context, urlStr string, form url.Values, result interface{}) error {
	req, err := http.NewRequest("POST", a.makeURL(urlStr), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return a.Do(ctx, req, result)
}

// Get sends a GET request to the urlStr and marshals the response into result
func (a *API) Get(ctx context.Context, urlStr string, result interface{}) error {
	req, err := http.NewRequest("GET", a.makeURL(urlStr), nil)
	if err != nil {
		return err
	}
	return a.Do(ctx, req, result)
}

func (a *API) makeURL(urlStr string) string {
	return fmt.Sprintf("https://%s.memberclicks.net/%s", a.orgID, strings.TrimPrefix(urlStr, "/"))
}

// Do sends the http.Request and marshals the JSON response into result
func (a *API) Do(ctx context.Context, req *http.Request, result interface{}) error {

	// If no authorization header already set, prefer accessToken
	if a.accessToken != "" && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+a.accessToken)
	}

	// Otherwise, set basic auth.
	if req.Header.Get("Authorization") == "" {
		req.SetBasicAuth(a.clientID, a.clientSecret)
	}

	// Set other general headers.
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "application/json")

	client := a.getClient(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(string(bodyBytes))
	}

	return json.Unmarshal(bodyBytes, result)
}

// New creates a new API client
func New(orgID, clientID, clientSecret string) *API {
	return &API{orgID: orgID, clientID: clientID, clientSecret: clientSecret}
}
