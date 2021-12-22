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
		converted = Map(Wrap(items), conv.And(toString, addTail), func(v int) bool { return v%2 == 0 })
	)
	assert.Equal(t, slice.Of("2_tail", "4_tail"), SliceOf(converted))

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
		multiDimension    = Of(Of(Of(1, 2, 3), Of(4, 5, 6)), Of(Of(7), nil), nil)
		oneDimension      = Filter(Flatt(Flatt(multiDimension, conv.To[Iterator[Iterator[int]]]), conv.To[Iterator[int]]), odds)
	)

	e := slice.Of(1, 3, 5, 7)
	a := SliceOf(oneDimension)
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

// func Test_FlattDeepStructure(t *testing.T) {
// 	type (
// 		Attributes struct{ name string }
// 		Item       struct{ attributes []*Attributes }
// 	)

// 	var (
// 		items = New(&Item{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil)

// 		getName       = func(a *Attributes) string { return a.name }
// 		getAttributes = func(item *Item) []*Attributes { return item.attributes }
// 	)

// 	names := Map(NotNil(Flatt(NotNil(items), getAttributes)), getName)
// 	assert.Equal(t, Of("first", "second"), names)

// 	names = Map(Flatt(items, getAttributes, check.NotNil[*Item]), getName, check.NotNil[*Attributes])
// 	assert.Equal(t, Of("first", "second"), names)

// 	names = Flatt(items, func(item *Item) []string { return Map(item.attributes, getName, check.NotNil[*Attributes]) }, check.NotNil[*Item])
// 	assert.Equal(t, Of("first", "second"), names)

// 	//plain old style
// 	names = make([]string, 0)
// 	for _, i := range items {
// 		if i != nil {
// 			for _, a := range i.attributes {
// 				if a != nil {
// 					names = append(names, a.name)
// 				}
// 			}
// 		}
// 	}

// 	assert.Equal(t, Of("first", "second"), names)
// }
