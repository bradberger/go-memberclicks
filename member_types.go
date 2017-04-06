package memberclicks

import "golang.org/x/net/context"

// MemberTypes is a MemberType list
type MemberTypes []MemberType

// MemberType is a member type
type MemberType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type memberTypeResp struct {
	TotalCount  int         `json:"totalCount"`
	MemberTypes MemberTypes `json:"memberTypes"`
}

// MemberTypes returns a slice of member types for the account
func (a *API) MemberTypes(ctx context.Context) (MemberTypes, error) {
	var res memberTypeResp
	if err := a.Get(ctx, "/api/v1/member-type", &res); err != nil {
		return nil, err
	}
	return res.MemberTypes, nil
}
