package gilk_test

import (
	"testing"

	test "github.com/stretchr/testify/assert"
)

func TestGilk(t *testing.T) {
	assert := test.New(t)

	// assert equality
	assert.Equal(123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(123, 456, "they should not be equal")
}
