package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_ToMapResolvOrder(t *testing.T) {

	var grouped map[string]int
	var order []string

	order, grouped = slice.ToMapResolvOrder(users, User.Name, User.Age, resolv.First)

	assert.Equal(t, 4, len(grouped))
	assert.Equal(t, slice.Of("Bob", "Alice", "Tom", "Chris"), order)

}
