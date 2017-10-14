package memberclicks

import (
	"fmt"

	"golang.org/x/net/context"
)

type ProfileResp struct {
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

// Profiles returns a page of profiles. Set the pageNum to be < 1 to get all pages at the same time.
func (a *API) Profiles(ctx context.Context, pageNum, pageSize int) (*ProfileResp, error) {

	all := pageNum < 1
	if all {
		pageSize = 100
		pageNum = 1
	}

	var resp ProfileResp
	if err := a.Get(ctx, fmt.Sprintf("/api/v1/profile?pageNumber=%d&pageSize=%d", pageNum, getPageSize(pageSize)), &resp); err != nil {
		return nil, err
	}

	if all {
		for i := 1; i < resp.TotalPageCount; i++ {
			pg, err := a.Profiles(ctx, i+1, pageSize)
			if err != nil {
				return &resp, err
			}
			resp.Profiles = append(resp.Profiles, pg.Profiles...)
		}
		resp.TotalPageCount = 1
		resp.TotalCount = len(resp.Profiles)
		resp.Count = resp.TotalCount
		resp.NextPageURL = ""
		resp.LastPageURL = ""
	}

	return &resp, nil
}

// ProfileSearch returns a profile search with the given ID. If pageNum is less than 1, all pages of the search will be returned.
func (a *API) ProfileSearch(ctx context.Context, searchID string, pageNum int) (*ProfileResp, error) {

	all := pageNum < 1
	if all || pageNum < 1 {
		pageNum = 1
	}

	var resp ProfileResp
	urlStr := fmt.Sprintf("/api/v1/profile?searchId=%s&pageSize=100&pageNumber=%d", searchID, pageNum)
	if err := a.Get(ctx, urlStr, &resp); err != nil {
		return nil, err
	}

	if all && resp.TotalPageCount > 1 {
		for i := 1; i < resp.TotalPageCount; i++ {
			pg, err := a.ProfileSearch(ctx, searchID, i+1)
			if err != nil {
				return &resp, err
			}
			resp.Profiles = append(resp.Profiles, pg.Profiles...)
		}
	}

	return &resp, nil
}

// ProfileSearch returns a profile search with the given ID.
func (a *API) GetProfileSearch(ctx context.Context, search *ProfileSearchResp) (*ProfileResp, error) {
	var resp ProfileResp
	if err := a.Get(ctx, search.ProfilesURL, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
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
	var resp ProfileResp
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
