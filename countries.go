package memberclicks

import "golang.org/x/net/context"

// Countries is a list of countries
type Countries []Country

// Country is a country struct
type Country struct {
	Name string `json:"name"`
}

func (c Country) String() string {
	return c.Name
}

type countryResponse struct {
	TotalCount int       `json:"totalCount"`
	Countries  Countries `json:"countries"`
}

// Countries returns a list of countries
func (a *API) Countries(ctx context.Context) (Countries, error) {
	var res countryResponse
	if err := a.Get(ctx, "/api/v1/country", &res); err != nil {
		return nil, err
	}
	return res.Countries, nil
}
