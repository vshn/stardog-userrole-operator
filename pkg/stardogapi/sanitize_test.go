package stardogapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizePathValue(t *testing.T) {
	value := "/..%2Fadmin"
	expectedValue := "..admin"

	assert.Equal(t, expectedValue, sanitizePathValue(value))
}
