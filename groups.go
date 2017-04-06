package memberclicks

import "golang.org/x/net/context"

// Groups is an slice of groups
type Groups []Group

// Group is a group
type Group struct {
	Name string `json:"name"`
}

type groupResp struct {
	TotalCount int    `json:"totalCount"`
	Groups     Groups `json:"groups"`
}

// Groups returns a list of groups for the acccount
func (a *API) Groups(ctx context.Context) (Groups, error) {
	var res groupResp
	if err := a.Get(ctx, "/api/v1/group", &res); err != nil {
		return nil, err
	}
	return res.Groups, nil
}
