package memberclicks

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
)

var (
	// ErrNoContent is returned when an HTTP resposne has a 2XX response code but a 0 ContentLength
	ErrNoContent = errors.New("HTTP response has no content length")
)

// NewAPI returns a new API instance with the given authToken
func NewAPI(org string) *API {
	return &API{OrganizationID: org}
}

func NewAPIWithClientIDAndSecret(organization, clientID, clientSecret string) *API {
	return &API{
		OrganizationID: organization,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
	}
}

// API handles sending/receiving data to/from the MemberClicks API
type API struct {
	OrganizationID, ClientID, ClientSecret string
	AccessToken                            *AccessToken
	authToken, endpoint                    string

	// Unit testing private variables
	lastRequest  *http.Request
	lastResponse *http.Response
}

func (a *API) getEndpoint() string {
	return fmt.Sprintf("https://%s.memberclicks.net", a.OrganizationID)
}

func (a *API) makeURL(uri string) string {
	return fmt.Sprintf("%s/%s", a.getEndpoint(), strings.TrimPrefix(strings.TrimPrefix(uri, a.getEndpoint()), "/"))
}

// PostXML sends a HTTP Post request to the given url with the reqData encoded as XML in the body, and returns the result
// decoded into respData.
// func (a *API) PostXML(ctx context.Context, uri string, reqData interface{}, respData interface{}) (*http.Response, error) {
// 	var body bytes.Buffer
// 	if err := xml.NewEncoder(&body).Encode(reqData); err != nil {
// 		return nil, err
// 	}
// 	req, err := http.NewRequest("POST", a.makeURL(uri), &body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return a.do(ctx, req, respData)
// }

// Post sends a HTTP Post request to the given url and returns the result decoded into respData.
func (a *API) Post(ctx context.Context, uri string, form url.Values, respData interface{}) (*http.Response, error) {
	req, err := http.NewRequest("POST", a.makeURL(uri), bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	return a.do(ctx, req, respData)
}

// Get sends a HTTP GET request to the given uri and returns the XML encoded respone in respData
func (a *API) Get(ctx context.Context, uri string, respData interface{}) (*http.Response, error) {
	req, err := http.NewRequest("GET", a.makeURL(uri), nil)
	if err != nil {
		return nil, err
	}
	return a.do(ctx, req, respData)
}

func (a *API) do(ctx context.Context, req *http.Request, respData interface{}) (*http.Response, error) {

	// For unit testing
	a.lastRequest = req

	switch {
	case a.AccessToken != nil:
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.AccessToken.AccessToken))
	case a.authToken != "":
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.authToken))
	case a.ClientID != "" && a.ClientSecret != "":
		req.SetBasicAuth(a.ClientID, a.ClientSecret)
	}

	var err error
	client := a.getClient(ctx)
	resp, err := client.Do(req)
	a.lastResponse = resp
	if err != nil {
		return resp, err
	}
	// TODO Not sure how to handle errors, what format they're in, etc.
	if resp.StatusCode >= 400 {
		// Try to format into an error response and get the message.
		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
			var errResp ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
				return resp, errors.New(errResp.Message)
			}
		}
		// Otherwise then just return the whole thing.
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		defer resp.Body.Close()
		return resp, fmt.Errorf("%s", string(bodyBytes))
	}

	switch {
	case strings.HasPrefix(resp.Header.Get("Content-Type"), "image/jp"):
		// TODO
	case strings.HasPrefix(resp.Header.Get("Content-Type"), "image/png"):
		// TODO
	case strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json"):
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

type ErrorResponse struct {
	Timestamp      int64             `json:"timestamp"`
	Status         int               `json:"status"`
	Error          string            `json:"error"`
	Message        string            `json:"message"`
	MessageDetails []string          `json:"messageDetails"` // not sure what format this is
	Path           string            `json:"path"`
	Parameters     map[string]string // not sure what format this is, probably string
}

func (e *ErrorResponse) Time() time.Time {
	return time.Unix(e.Timestamp, 0)
}
