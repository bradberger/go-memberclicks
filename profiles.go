package memberclicks

import (
	"fmt"

	"golang.org/x/net/context"
)

type profileResp struct {
	TotalCount     int       `json:"totalCount"`
	TotalPageCount int       `json:"totalPageCount"`
	PageNumber     int       `json:"pageNumber"`
	PageSize       int       `json:"pageSize"`
	Count          int       `json:"count"`
	FirstPageURL   string    `json:"firstPageUrl"`
	NextPageURL    string    `json:"nextPageUrl"`
	LastPageURL    string    `json:"lastPageUrl"`
	Profiles       []Profile `json:"profiles"`
}

type ProfileSearchResp struct {
	ID          string                 `json:"id"`
	ExpireDate  string                 `json:"expireDate"`
	Status      int                    `json:"status"`
	Timestamp   int64                  `json:"timestamp"`
	URL         string                 `json:"url"`
	Item        map[string]interface{} `json:"item"`
	ProfilesURL string                 `json:"profilesUrl"`
}

// Profiles returns a page of profiles
func (a *API) Profiles(ctx context.Context, page, pageSize int) ([]Profile, error) {

	var resp profileResp
	if page < 1 {
		page = 1
	}

	if err := a.Get(ctx, fmt.Sprintf("/api/v1/profile?pageNumber=%d&pageSize=%d", page, getPageSize(pageSize)), &resp); err != nil {
		return nil, err
	}

	return resp.Profiles, nil
}

// ProfileSearch returns a profile search with the given ID.
func (a *API) ProfileSearch(ctx context.Context, searchID string) ([]Profile, error) {
	var resp profileResp
	urlStr := "/api/v1/profile/search/" + searchID + "?pageSize=100"
	if true {
		urlStr = "/api/v1/profile?searchId=" + searchID + "&pageSize=100"
	}
	if err := a.Get(ctx, urlStr, &resp); err != nil {
		return nil, err
	}
	return resp.Profiles, nil
}

// ProfileSearch returns a profile search with the given ID.
func (a *API) GetProfileSearch(ctx context.Context, search *ProfileSearchResp) ([]Profile, error) {
	var resp profileResp
	if err := a.Get(ctx, search.ProfilesURL, &resp); err != nil {
		return nil, err
	}
	return resp.Profiles, nil
}

// CreateProfileSearch creates a profile search and returns the resulting search ID.
// Profile searches exipre every half an hour.
func (a *API) CreateProfileSearch(ctx context.Context, params interface{}) (*ProfileSearchResp, error) {
	var resp ProfileSearchResp
	if err := a.PostJSON(ctx, "/api/v1/profile/search", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProfilePageCt gets the total number of profiles
func (a *API) ProfilePageCt(ctx context.Context, pageSize int) (int, error) {
	if pageSize < 10 {
		pageSize = 10
	}
	var resp profileResp
	if err := a.Get(ctx, fmt.Sprintf("/api/v1/profile?pageNumber=1&pageSize=%d", getPageSize(pageSize)), &resp); err != nil {
		return 0, err
	}
	return resp.TotalPageCount, nil
}

func getPageSize(pageSize int) int {
	switch {
	case pageSize < 10:
		return 10
	case pageSize > 100:
		return 100
	default:
		return pageSize
	}
}
