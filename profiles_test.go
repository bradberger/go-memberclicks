package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfiles(t *testing.T) {
	p, err := mc.Profiles(ctx, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, p, 10)
	assert.True(t, p[0].ID() > 0)
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
