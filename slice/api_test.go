package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Range(t *testing.T) {
	assert.Equal(t, Of(-1, 0, 1, 2, 3), Range(-1, 3))
	assert.Equal(t, Of(3, 2, 1, 0, -1), Range(3, -1))
	assert.Equal(t, Of(1), Range(1, 1))
}
