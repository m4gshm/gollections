package group

import (
	"testing"

	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/iter"

	"github.com/stretchr/testify/assert"
)

func Test_group_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = Of(ptr.Of(iter.New(slice.Of(1, 1, 2, 4, 3, 1))).Next, even)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}
