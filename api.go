package memberclicks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// Scopes
const (
	ScopeRead = "read"
)

var (
	// Timeout is the duration before HTTP requests to the API servers time out
	Timeout = 15 * time.Second
)

// API is an MemberClicks API client
type API struct {
	orgID, clientID, clientSecret, accessToken string

	Client  *http.Client
	Timeout time.Duration
}

func (a *API) getTimeout() time.Duration {
	if a.Timeout > 0 {
		return a.Timeout
	}
	return Timeout
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

func (a *API) GetAuthRequestURL(scope, state, redirectURL string) string {
	return fmt.Sprintf(
		"https://%s.memberclicks.net/oauth/v1/authorize?response_type=token&client_id=%s&scope=%s&state=%s&redirect_uri=%s",
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
	if err := a.Post(ctx, "/oauth/v1/token", form, &t); err != nil {
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

// OwnerPassword returns a token with the owner "password" grant type
func (a *API) OwnerPassword(ctx context.Context, username, password string) (*Token, error) {
	var t Token
	form := url.Values{
		"scope":      {"read"},
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	if err := a.Post(ctx, "/oauth/v1/token", form, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// CheckPassword is shorthand for OwnerPasswordGrant, without returning the token
func (a *API) CheckPassword(ctx context.Context, username, password string) (err error) {
	_, err = a.OwnerPassword(ctx, username, password)
	return
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

func (a *API) ResourceOwnerCredentials(ctx context.Context, username, password string) (*Token, error) {
	var t Token
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("scope", "read")
	form.Add("username", username)
	form.Add("password", password)
	if err := a.Post(ctx, "/oauth/v1/token", form, &t); err != nil {
		return nil, err
	}
	return &t, nil
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

// PostJSON sends a JSON POST request to the API with the JSON encoded data
func (a *API) PostJSON(ctx context.Context, urlStr string, data, result interface{}) error {

	buf := bytes.NewBuffer(nil)
	buf2 := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(io.MultiWriter(buf, buf2)); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", a.makeURL(urlStr), buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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

	if resp.StatusCode >= 400 {
		return errors.New(string(bodyBytes))
	}

	return json.Unmarshal(bodyBytes, result)
}

// New creates a new API client
func New(orgID, clientID, clientSecret string) *API {
	return &API{orgID: orgID, clientID: clientID, clientSecret: clientSecret}
}
