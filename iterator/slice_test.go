package iterator

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/slice"
	"github.com/stretchr/testify/assert"
)

func Test_MapAndFilter(t *testing.T) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }

		items     = slice.Of(1, 2, 3, 4, 5)
		converted = Map(Iter(items), conv.And(toString, addTail), func(v int) bool { return v%2 == 0 })
	)
	assert.Equal(t, slice.Of("2_tail", "4_tail"), ToSlice(converted))

	//plain old style
	convertedOld := make([]string, 0)
	for _, i := range items {
		if i%2 == 0 {
			convertedOld = append(convertedOld, conv.And(toString, addTail)(i))
		}
	}

	assert.Equal(t, slice.Of("2_tail", "4_tail"), convertedOld)

}

func Test_FlattSlices(t *testing.T) {
	var (
		odds              = func(v int) bool { return v%2 != 0 }
		multiDimensionOld = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
		multiDimension    = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}

		oneDimension      = Filter(Flatt(Flatt(Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds)
	)

	e := slice.Of(1, 3, 5, 7)
	a := ToSlice(oneDimension)
	assert.Equal(t, e, a)

	//plain old style
	oneDimensionOld := make([]int, 0)
	for _, i := range multiDimensionOld {
		if i == nil {
			continue
		}
		for _, ii := range i {
			if ii == nil {
				continue
			}
			for _, iii := range ii {
				if odds(iii) {
					oneDimensionOld = append(oneDimensionOld, iii)
				}
			}
		}
	}

	assert.Equal(t, slice.Of(1, 3, 5, 7), oneDimensionOld)

}