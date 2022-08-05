package deque_test

import (
	"testing"

	test "github.com/stretchr/testify/assert"
)

// TODO: add test
func TestDeque(t *testing.T) {
	t.Parallel()

	assert := test.New(t)

	// assert equality
	assert.Equal(1, 1, "they should be equal")

	// assert inequality
	assert.NotEqual(1, 2, "they should not be equal")
}
