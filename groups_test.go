package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroups(t *testing.T) {
	g, err := mc.Groups(ctx)
	assert.NoError(t, err)
	assert.True(t, len(g) > 1)
}
