package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	breakgroup "github.com/m4gshm/gollections/break/kv/loop/group"
	breakloop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/kv/loop/group"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/iter"
)

func Test_group_odd_even(t *testing.T) {
	var (
		even     = func(v int) bool { return v%2 == 0 }
		elements = slice.Of(1, 1, 2, 4, 3, 1)
		it       = iter.ToKV(elements, even, as.Is[int])
		groups   = group.Of(it.Next)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_group_odd_even2(t *testing.T) {
	var (
		even     = func(v int) bool { return v%2 == 0 }
		elements = slice.Of(1, 1, 2, 4, 3, 1)
		it       = iter.ToKVs(elements, func(i int) []bool { return slice.Of(even(i)) }, func(i int) []int { return slice.Of(i) })
		groups   = group.Of(it.Next)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_group_odd_even3(t *testing.T) {
	var (
		even      = func(v int) bool { return v%2 == 0 }
		elements  = slice.Of(1, 1, 2, 4, 3, 1)
		it        = iter.NewKeyVal(elements, func(i int) (bool, error) { return even(i), nil }, func(i int) (int, error) { return i, nil })
		groups, _ = breakgroup.Of(it.Next)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_group_odd_even4(t *testing.T) {
	var (
		even      = func(v int) bool { return v%2 == 0 }
		elements  = slice.Of(1, 1, 2, 4, 3, 1)
		it        = iter.NewMultipleKeyVal(elements, func(i int) ([]bool, error) { return slice.Of(even(i)), nil }, func(i int) ([]int, error) { return slice.Of(i), nil })
		groups, _ = breakgroup.Of(it.Next)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}

func Test_Filter(t *testing.T) {
	var (
		even     = func(v int) bool { return v%2 == 0 }
		elements = slice.Of(1, 1, 2, 4, 6, 3, 1)
		it       = iter.Filter(elements, even)
		out      = loop.Slice(it.Next)
	)
	assert.Equal(t, slice.Of(2, 4, 6), out)
}

func Test_Filt(t *testing.T) {
	var (
		even     = func(v int) bool { return v%2 == 0 }
		elements = slice.Of(1, 1, 2, 4, 6, 3, 1)
		it       = iter.Filt(elements, func(i int) (bool, error) { return even(i), nil })
		out, _   = breakloop.Slice(it.Next)
	)
	assert.Equal(t, slice.Of(2, 4, 6), out)
}

func Test_Convert(t *testing.T) {
	var (
		elements = slice.Of(1, 2, 4, 6)
		it       = iter.Convert(elements, strconv.Itoa)
		out      = loop.Slice(it.Next)
	)
	assert.Equal(t, slice.Of("1", "2", "4", "6"), out)
}

func Test_Conv(t *testing.T) {
	var (
		elements = slice.Of("1", "2", "4", "6")
		it       = iter.Conv(elements, strconv.Atoi)
		out, _   = breakloop.Slice(it.Next)
	)
	assert.Equal(t, slice.Of(1, 2, 4, 6), out)
}
