package memberclicks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	var a, b string
	var c int
	a = "foo"

	assert.Error(t, copy(a, b))
	assert.Error(t, copy(a, &c))
	assert.NoError(t, copy(&a, &b))
	assert.Equal(t, a, b)
}
