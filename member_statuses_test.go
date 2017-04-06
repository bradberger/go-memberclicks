package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemberStatusesAPI(t *testing.T) {
	s, err := mc.MemberStatuses(ctx)
	assert.NoError(t, err)
	assert.True(t, len(s) > 0)
}
