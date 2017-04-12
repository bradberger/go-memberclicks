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
