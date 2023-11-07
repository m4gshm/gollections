package boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/m4gshm/gollections/slice/first"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/m4gshm/gollections/slice/reverse"
)

func Test_Slice_schortcuts(t *testing.T) {

	data := slice.Of("Bob", "Chris", "Alice") // constructor

	sorted := sort.Asc(data) //sorting

	reversed := reverse.Of(clone.Of(sorted)) //reversing of cloned slice

	chris, ok := first.Of(reversed, func(name string) bool { return name[0] == 'C' }) //finding the first element by a predicate function

	var lengthMap map[int][]string = group.Of(sorted, func(name string) int { return len(name) }, as.Is[string]) // converting to a map

	assert.Equal(t, slice.Of("Alice", "Bob", "Chris"), sorted)
	assert.Equal(t, "Chris", chris)
	assert.True(t, true, ok)
	assert.Equal(t, slice.Of("Alice", "Chris"), lengthMap[5])

}
