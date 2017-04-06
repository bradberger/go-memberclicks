package memberclicks

// Token is a OAuth2 access token response from the server
type Token struct {
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
