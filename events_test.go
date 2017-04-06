package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventsAPI(t *testing.T) {
	e, err := mc.Events(ctx)
	assert.NoError(t, err)
	assert.True(t, len(e) > 0)
}
