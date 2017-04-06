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

// Profiles returns a page of profiles
func (a *API) Profiles(ctx context.Context, page int) ([]Profile, error) {

	var resp profileResp
	if page < 1 {
		page = 1
	}

	if err := a.Get(ctx, fmt.Sprintf("/api/v1/profile?pageNumber=%d", page), &resp); err != nil {
		return nil, err
	}

	return resp.Profiles, nil
}

// ProfilePageCt gets the total number of profiles
func (a *API) ProfilePageCt(ctx context.Context) (int, error) {
	var resp profileResp
	if err := a.Get(ctx, "/api/v1/profile?pageNumber=1", &resp); err != nil {
		return 0, err
	}
	return resp.TotalPageCount, nil
}
