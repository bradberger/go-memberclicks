package memberclicks

import (
	"net/url"

	"golang.org/x/net/context"
)

const (
	// ScopeRead is the read scope for access tokens requests
	ScopeRead = "read"

	// GrantTypeClientCredentials is the client credentials grant type
	GrantTypeClientCredentials = "client_credentials"
)

// AccessToken is a OAuth2 access token response from the server
type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	// ExpiresIn is expiration time in seconds, usually 3600 or one hour
	ExpiresIn int64  `json:"expires_in"`
	Scope     string `json:"scope"`
	ServiceID int64  `json:"serviceId"`
	UserID    int64  `json:"userId"`
	JTI       string `json:"jti"`
}

// GetAccessToken gets client credentials from the API with a client credentials grant type.
func (a *API) GetAccessToken(ctx context.Context, scope string) (*AccessToken, error) {
	var t AccessToken
	if scope == "" {
		scope = ScopeRead
	}
	if _, err := a.Post(ctx, "/oauth/v1/token", url.Values{"grant_type": {GrantTypeClientCredentials}, "scope": {scope}}, &t); err != nil {
		return nil, err
	}
	a.AccessToken = &t
	return a.AccessToken, nil
}
