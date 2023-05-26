package group

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/break/loop"
)

func Test_group_odd_even(t *testing.T) {

	var (
		even      = func(v int) (bool, error) { return v%2 == 0, nil }
		groups, _ = Of(loop.NewKeyValuer(loop.Of(1, 1, 2, 4, 3, 1), even, func(i int) (int, error) { return i, nil }).Next)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}
