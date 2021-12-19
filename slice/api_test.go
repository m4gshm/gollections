package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MapAndFilter(t *testing.T) {
	var toString Converter[int, string] = func(i int) string { return fmt.Sprintf("%d", i) }
	var addTail Converter[string, string] = func(s string) string { return s + "_tail" }
	var evens Predicate[int] = func(v int) bool { return v%2 == 0 }

	converted := Map(Of(1, 2, 3, 4, 5), And(toString, addTail), evens)

	assert.Equal(t, Of("2_tail", "4_tail"), converted)
}

func Test_FlattSlices(t *testing.T) {
	var (
		multiDimension [][][]int
		oneDimension   []int
	)

	var odds Predicate[int] = func(v int) bool { return v%2 != 0 }

	multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	oneDimension = Flatt(Flatt(multiDimension, AsIs[[][]int]), AsIs[[]int], odds)

	assert.Equal(t, Of(1, 3, 5, 7), oneDimension)
}

func Test_FlattDeepStructure(t *testing.T) {
	type (
		Attributes struct {
			name string
		}
		Item struct {
			attributes []*Attributes
		}
	)
	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)

	items := []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}}

	names := Map(Flatt(items, getAttributes, NotNil[*Attributes]), getName)

	assert.Equal(t, Of("first", "second"), names)
}
