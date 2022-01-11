package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/conv"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/op"
	slc "github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_OrderedSet(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1)
	values := s.Collect()
	fmt.Println(s) //[1, 2, 4, 3]

	assert.Equal(t, slc.Of(1, 2, 4, 3), values)
}

func Test_group_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.Of(set.Of(1, 1, 2, 4, 3, 1), even)
	)
	fmt.Println(groups) //map[false:[1 3] true:[2 4]]
	assert.Equal(t, map[bool][]int{false: {1, 3}, true: {2, 4}}, groups)
}

func Test_group_with_filtering_by_stirng_len(t *testing.T) {
	var groups = it.Group(
		set.Of(
			"seventh",
			"first", "second", "third", "fourth",
			"fifth", "sixth", "eighth",
			"ninth", "tenth", "one", "two", "three", "1",
		).Begin(),
		func(v string) int { return len(v) },
	).Map(func(kv *typ.KV[int, string]) *typ.KV[int, string] {
		val := kv.Value()
		if val == "sixth" {
			return nil
		}
		return K.V(kv.Key(), val+"_")
	}).Filter(
		check.NotNil[typ.KV[int, string]],
	).Filter(
		func(kv *typ.KV[int, string]) bool { return kv.Key() > 3 },
	).Collect()

	fmt.Println(groups) //map[5:[first_ third_ fifth_ ninth_ tenth_ three_] 6:[second_ fourth_ eighth_] 7:[seventh_]]

	assert.Equal(t, map[int][]string{
		5: {"first_", "third_", "fifth_", "ninth_", "tenth_", "three_"},
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
	oddSum := it.Reduce(it.Filter(it.Flatt(slc.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
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
