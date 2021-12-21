package slice

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/conv"
	"github.com/stretchr/testify/assert"
)

func Test_MapAndFilter(t *testing.T) {
	var (
		toString conv.Converter[int, string]    = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  conv.Converter[string, string] = func(s string) string { return s + "_tail" }

		items     = Of(1, 2, 3, 4, 5)
		converted = Map(items, conv.And(toString, addTail), func(v int) bool { return v%2 == 0 })
	)
	assert.Equal(t, Of("2_tail", "4_tail"), converted)

	//plain old style
	converted = make([]string, 0)
	for _, i := range items {
		if i%2 == 0 {
			converted = append(converted, conv.And(toString, addTail)(i))
		}
	}

	assert.Equal(t, Of("2_tail", "4_tail"), converted)

}

func Test_FlattSlices(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
		oneDimension   = Filter(Flatt(Flatt(NotNil(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds)
	)

	assert.Equal(t, Of(1, 3, 5, 7), oneDimension)

	//plain old style
	oneDimension = make([]int, 0)
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
					oneDimension = append(oneDimension, iii)
				}
			}
		}
	}

	assert.Equal(t, Of(1, 3, 5, 7), oneDimension)

}

func Test_FlattDeepStructure(t *testing.T) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
		names         = Map(NotNil(Flatt(NotNil(items), getAttributes)), getName)
	)

	assert.Equal(t, Of("first", "second"), names)

	//plain old style
	names = make([]string, 0)
	for _, i := range items {
		if i != nil {
			for _, a := range i.attributes {
				if a != nil {
					names = append(names, a.name)
				}
			}
		}
	}

	assert.Equal(t, Of("first", "second"), names)
}
