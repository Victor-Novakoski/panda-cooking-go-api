package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUUID(t *testing.T) {
	assert.True(t, isUUID("4a7f5c1e-2b3d-4c5e-8f9a-0b1c2d3e4f5a"))
	assert.True(t, isUUID("4A7F5C1E-2B3D-4C5E-8F9A-0B1C2D3E4F5A"))
	assert.False(t, isUUID(""))
	assert.False(t, isUUID("123"))
	assert.False(t, isUUID("4a7f5c1e-2b3d-4c5e-8f9a-0b1c2d3e4f5a' OR '1'='1"))
}
