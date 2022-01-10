package test

import (
	"testing"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/stretchr/testify/assert"
)

func Test_Range(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), range_.Of(-1, 3))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), range_.Of(3, -1))
	assert.Equal(t, slice.Of(1), range_.Of(1, 1))
}
