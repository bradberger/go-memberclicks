package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfiles(t *testing.T) {
	resp, err := mc.Profiles(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, resp.Profiles, 10)
	assert.True(t, resp.Profiles[0].ID() > 0)
}

func TestProfilesErr(t *testing.T) {
	_, err := mc.Profiles(ctx, 1000000, 10)
	assert.Error(t, err)
}

func TestProfilePageCt(t *testing.T) {
	ct, err := mc.ProfilePageCt(ctx, 0)
	assert.NoError(t, err)
	assert.True(t, ct > 0)
}

func TestGetPageSize(t *testing.T) {
	assert.Equal(t, 10, getPageSize(0))
	assert.Equal(t, 10, getPageSize(10))
	assert.Equal(t, 100, getPageSize(1000))
}

func TestProfileSearch(t *testing.T) {

	params := map[string]interface{}{"[Last Modified Date]": "Last 15 Minutes"}
	search, err := mc.CreateProfileSearch(ctx, &params)
	assert.NoError(t, err)
	if !assert.NotNil(t, search) {
		return
	}
	assert.NotEmpty(t, search.ID)

	_, err = mc.ProfileSearch(ctx, search.ID, 1)
	assert.NoError(t, err)
}
