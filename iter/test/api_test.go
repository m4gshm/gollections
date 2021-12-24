package iter

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/stretchr/testify/assert"
)

func Test_MapAndFilter(t *testing.T) {

	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
	)
	items := []int{1, 2, 3, 4, 5}
	converted := iter.MapFit(iter.Wrap(items), func(v int) bool { return v%2 == 0 }, conv.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), iter.ToSlice(converted))

	converted = slice.MapFit(items, func(v int) bool { return v%2 == 0 }, conv.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), iter.ToSlice(converted))

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
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	e := slice.Of(1, 3, 5, 7)

	a := iter.ToSlice(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds))
	assert.Equal(t, e, a)

	a = iter.ToSlice(iter.Filter(iter.Flatt(slice.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds))
	assert.Equal(t, e, a)

	//plain old style
	oneDimensionOld := make([]int, 0)
	for _, i := range multiDimension {
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

func Test_ReduceSlices(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	e := 1 + 3 + 5 + 7

	oddSum := iter.Reduce(iter.Filter(iter.Flatt(iter.FlattFit(iter.Wrap(multiDimension), check.NotNil[[][]int], conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
	assert.Equal(t, e, oddSum)

	oddSum = iter.Reduce(iter.Filter(iter.Flatt(slice.FlattFit(multiDimension, check.NotNil[[][]int], conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
	assert.Equal(t, e, oddSum)

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

	assert.Equal(t, e, oddSum)

}