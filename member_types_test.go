package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemberTypes(t *testing.T) {
	m, err := mc.MemberTypes(ctx)
	assert.NoError(t, err)
	assert.True(t, len(m) > 0)
}
