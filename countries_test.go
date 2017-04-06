package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountryStringer(t *testing.T) {
	c := Country{Name: "foo"}
	assert.Equal(t, "foo", c.String())
}

func TestCountriesAPI(t *testing.T) {
	c, err := mc.Countries(ctx)
	assert.NoError(t, err)
	assert.True(t, len(c) > 0)
}
