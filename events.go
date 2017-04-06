package memberclicks

import "golang.org/x/net/context"

// Events is a list of events
type Events []Event

// Event is an event
type Event struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type eventsResponse struct {
	TotalCount int    `json:"totalCount"`
	Events     Events `json:"events"`
}

// Events returns an event list
func (a *API) Events(ctx context.Context) (Events, error) {
	var res eventsResponse
	if err := a.Get(ctx, "/api/v1/event", &res); err != nil {
		return nil, err
	}
	return res.Events, nil
}
