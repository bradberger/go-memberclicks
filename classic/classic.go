package classic

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
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

// New returns a new Classic instance with the given authToken
func New(ctx context.Context, orgID, apiKey, username, password string) (*Client, error) {
	api := &Client{
		OrganizationID: orgID,
		apiKey:         apiKey,
		username:       username,
		password:       password,
	}

	return api, api.Auth(ctx)
}

// Classic handles sending/receiving data to/from the MemberClicks Classic
type Client struct {
	apiKey, username, password, token string
	OrganizationID                    string

	HttpClient *http.Client

	sync.Mutex
}

func (c *Client) Auth(ctx context.Context) error {
	var auth AuthResponse
	form := url.Values{"apiKey": {c.apiKey}, "username": {c.username}, "password": {c.password}}
	if _, err := c.Post(ctx, "/services/auth", form, &auth); err != nil {
		return err
	}
	c.token = auth.Token
	return nil
}

func (c *Client) getEndpoint() string {
	return fmt.Sprintf("https://%s.memberclicks.net", c.OrganizationID)
}

func (c *Client) makeURL(uri string) string {
	return fmt.Sprintf("%s/%s", c.getEndpoint(), strings.TrimPrefix(strings.TrimPrefix(uri, c.getEndpoint()), "/"))
}

// Post sends a HTTP Post request to the given url and returns the result decoded into respDatc.
func (c *Client) Post(ctx context.Context, uri string, form url.Values, respData interface{}) (*http.Response, error) {
	req, err := c.NewRequest("POST", c.makeURL(uri), bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, respData)
}

// PostXML sends a XML POST request to the urlStr with a body of data and decodes result into respData
func (c *Client) PostXML(ctx context.Context, urlStr string, data interface{}, respData interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", c.makeURL(urlStr), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
	return c.do(ctx, req, respData)
}

// PostJSON sends a JSON POST request to the urlStr with a body of data and decodes result into respData
func (c *Client) PostJSON(ctx context.Context, urlStr string, data interface{}, respData interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", c.makeURL(urlStr), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(ctx, req, respData)
}

// Get sends a HTTP GET request to the given uri and returns the XML encoded respone in respData
func (c *Client) Get(ctx context.Context, uri string, respData interface{}) (*http.Response, error) {
	req, err := c.NewRequest("GET", c.makeURL(uri), nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, respData)
}

// NewRequest matches the http.NewRequest signature but sets required headers.
func (c *Client) NewRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.makeURL(urlStr), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

func (c *Client) Users(ctx context.Context, params url.Values) (*UserList, error) {
	var l UserList
	urlStr := "/services/user"
	if params != nil {
		urlStr += "?" + params.Encode()
	}
	if _, err := c.Get(ctx, urlStr, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

func (c *Client) User(ctx context.Context, userID string) (*Profile, error) {
	var p Profile
	if _, err := c.Get(ctx, "/services/user/"+userID, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, respData interface{}) (*http.Response, error) {

	var err error
	client := c.getClient(ctx)
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

	if true {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		fmt.Printf("%s\n", bodyBytes)
	}

	switch resp.Header.Get("Content-Type") {
	case "application/json":
		if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return resp, err
		}
	default:
		if err := xml.NewDecoder(resp.Body).Decode(respData); err != nil {
			return resp, err
		}
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
