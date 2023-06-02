package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/loop/group"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_group_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.Of(slice.NewIter(slice.Of(1, 1, 2, 4, 3, 1)).Next, even, as.Is[int])
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_ByMultiple(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.ByMultiple(slice.NewIter(slice.Of(1, 1, 2, 4, 3, 1)).Next, func(i int) []bool { return slice.Of(even(i)) }, as.Slice[int])
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_ByMultipleEmptyKey(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.ByMultiple(slice.NewIter(slice.Of(1, 1, 2, 4, 3, 1)).Next, func(i int) []bool { return op.IfElse(even(i), slice.Of(true), nil) }, as.Slice[int])
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_ByMultipleEmptyVal(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.ByMultiple(slice.NewIter(slice.Of(1, 1, 2, 4, 3, 1)).Next,
			func(i int) []bool { return slice.Of(even(i)) },
			func(i int) []int { return op.IfElse(even(i), nil, slice.Of(i)) },
		)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {0, 0}}, groups)
}

func Test_ByMultipleKeys(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.ByMultipleKeys(slice.NewIter(slice.Of(1, 1, 2, 4, 3, 1)).Next, func(i int) []bool { return slice.Of(even(i)) }, as.Is[int])
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_ByMultipleValues(t *testing.T) {

	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.ByMultipleValues(slice.NewIter(slice.Of(1, 1, 2, 4, 3, 1)).Next, even, as.Slice[int])
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}
