package test

import (
	"fmt"
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

func Test_StringRepresentation(t *testing.T) {
	var (
		expected = fmt.Sprint(slice.Of(1, 2, 3, 4))
		actual   = slice.ToString(slice.Of(1, 2, 3, 4))
	)
	assert.Equal(t, expected, actual)
}

func Test_StringReferencesRepresentation(t *testing.T) {
	var (
		expected       = fmt.Sprint(slice.Of(1, 2, 3, 4))
		i1, i2, i3, i4 = 1, 2, 3, 4
		actual         = slice.ToStringRefs(slice.Of(&i1, &i2, &i3, &i4))
	)
	assert.Equal(t, expected, actual)
}
