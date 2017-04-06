package memberclicks

import "golang.org/x/net/context"

// MemberStatuses is a complete list of member statuses.
type MemberStatuses []MemberStatus

// MemberStatus is a member status
type MemberStatus struct {
	Name string `json:"name"`
}

type memberStatusResp struct {
	TotalCount     int            `json:"totalCount"`
	MemberStatuses MemberStatuses `json:"memberStatuses"`
}

// MemberStatuses returns the complete list of member statuses.
func (a *API) MemberStatuses(ctx context.Context) (MemberStatuses, error) {
	var res memberStatusResp
	if err := a.Get(ctx, "/api/v1/member-status", &res); err != nil {
		return nil, err
	}
	return res.MemberStatuses, nil
}
