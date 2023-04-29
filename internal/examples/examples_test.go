package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/as"
	cGroup "github.com/m4gshm/gollections/collection/group"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Set(t *testing.T) {
	var (
		s      immutable.Set[int] = set.Of(1, 1, 2, 4, 3, 1)
		values []int              = s.Slice()
	)

	assert.Equal(t, 4, s.Len())
	assert.Equal(t, 4, len(values))

	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
	assert.True(t, s.Contains(4))
	assert.False(t, s.Contains(5))
}

func Test_OrderedSet(t *testing.T) {
	s := oset.Of(1, 1, 2, 4, 3, 1)
	values := s.Slice()
	fmt.Println(s) //[1, 2, 4, 3]

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_group_orderset_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.Of(oset.Of(1, 1, 2, 4, 3, 1), even)
	)
	fmt.Println(groups) //map[false:[1 3] true:[2 4]]
	assert.Equal(t, map[bool][]int{false: {1, 3}, true: {2, 4}}, groups)
}

func Test_group_orderset_with_filtering_by_stirng_len(t *testing.T) {
	var groups = cGroup.Of(oset.Of(
		"seventh", "seventh", //duplicated
		"first", "second", "third", "fourth",
		"fifth", "sixth", "eighth",
		"ninth", "tenth", "one", "two", "three", "1",
		"second", //duplicate
	), func(v string) int { return len(v) },
	).FilterKey(
		more.Than(3),
	).ConvertValue(
		func(v string) string { return v + "_" },
	).Map()

	assert.Equal(t, map[int][]string{
		5: {"first_", "third_", "fifth_", "sixth_", "ninth_", "tenth_", "three_"},
		6: {"second_", "fourth_", "eighth_"},
		7: {"seventh_"},
	}, groups)
}

func Test_compute_odds_sum(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
		expected       = 1 + 3 + 5 + 7
	)

	//declarative style
	oddSum := iter.Reduce(iter.Filter(iter.Flatt(sliceIter.Flatt(multiDimension, as.Is[[][]int]), as.Is[[]int]), odds), op.Sum[int])
	assert.Equal(t, expected, oddSum)

	//plain old style
	oddSum = 0
	for _, i := range multiDimension {
		for _, ii := range i {
			for _, iii := range ii {
				if odds(iii) {
					oddSum += iii
				}
			}
		}
	}

	assert.Equal(t, expected, oddSum)
}
