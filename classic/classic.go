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
)

var (
	// ErrNoContent is returned when an HTTP resposne has a 2XX response code but a 0 ContentLength
	ErrNoContent = errors.New("HTTP response has no content length")
)

// AuthResponse is a MemberClicks classic auth response
type AuthResponse struct {
	UserID        string `json:"userId"`
	GroupID       string `json:"groupId"`
	OrgID         string `json:"orgId"`
	ContactName   string `json:"contactName"`
	UserName      string `json:"userName"`
	Active        bool   `json:"active,string"`
	Validated     bool   `json:"validated,string"`
	Deleted       bool   `json:"deleted,string"`
	FormStatus    string `json:"formStats"`
	LastModify    string `json:"lastModify"`
	NoMassEmail   bool   `json:"noMassEmail,string"`
	PrefBBContact string `json:"prefBBContact"`
	PrefBBImage   string `json:"prefBBImage"`
	Token         string `json:"token"`
	Password      string `json:"password"`
}

// NewClassic returns a new Classic instance with the given authToken
func NewClassic(orgID, apiKey, username, password string) (*Classic, error) {
	api := &Classic{
		OrganizationID: orgID,
		apiKey:         apiKey,
		username:       username,
		password:       password,
	}

	return api, api.getToken()
}

// Classic handles sending/receiving data to/from the MemberClicks Classic
type Classic struct {
	apiKey, username, password, token string
	OrganizationID                    string
}

func (a *Classic) getToken() error {
	var auth AuthResponse
	form := url.Values{"apiKey": {a.apiKey}, "username": {a.username}, "password": {a.password}}
	if _, err := a.Post("/services/auth", form, &auth); err != nil {
		return err
	}
	a.token = auth.Token
	return nil
}

func (a *Classic) getEndpoint() string {
	return fmt.Sprintf("https://%s.memberclicks.net", a.OrganizationID)
}

func (a *Classic) makeURL(uri string) string {
	return fmt.Sprintf("%s/%s", a.getEndpoint(), strings.TrimPrefix(strings.TrimPrefix(uri, a.getEndpoint()), "/"))
}

// Post sends a HTTP Post request to the given url and returns the result decoded into respData.
func (a *Classic) Post(uri string, form url.Values, respData interface{}) (*http.Response, error) {
	req, err := a.NewRequest("POST", a.makeURL(uri), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	return a.do(req, respData)
}

// Get sends a HTTP GET request to the given uri and returns the XML encoded respone in respData
func (a *Classic) Get(uri string, respData interface{}) (*http.Response, error) {
	req, err := a.NewRequest("GET", a.makeURL(uri), nil)
	if err != nil {
		return nil, err
	}
	return a.do(req, respData)
}

// NewRequest matches the http.NewRequest signature but sets required headers.
func (a *Classic) NewRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, a.makeURL(urlStr), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if a.token != "" {
		req.Header.Set("Authorization", a.token)
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

func (a *Classic) getClient() *http.Client {
	return httpClient
}

func (a *Classic) do(req *http.Request, respData interface{}) (*http.Response, error) {

	var err error
	client := a.getClient()
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	// If not 200, then return the status text.
	if resp.StatusCode >= 400 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		return resp, fmt.Errorf("error: %s", string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return resp, err
	}

	return resp, nil
}

// ErrorResponse is a MemberClicks classic error response
type ErrorResponse struct {
	Timestamp      int64             `json:"timestamp"`
	Status         int               `json:"status"`
	Error          string            `json:"error"`
	Message        string            `json:"message"`
	MessageDetails []string          `json:"messageDetails"` // not sure what format this is
	Path           string            `json:"path"`
	Parameters     map[string]string // not sure what format this is, probably string
}

// Time returns the time of the error response
func (e *ErrorResponse) Time() time.Time {
	return time.Unix(e.Timestamp, 0)
}
