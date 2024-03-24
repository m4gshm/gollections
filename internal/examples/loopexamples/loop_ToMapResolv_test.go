package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
)

func Test_Slice_ToMapResolv(t *testing.T) {

	var ageGroupedSortedNames map[string][]string

	ageGroupedSortedNames = loop.ToMapResolv(loop.Of(users...), func(u User) string {
		return op.IfElse(u.age <= 30, "<=30", ">30")
	}, User.Name, resolv.SortedSlice)

	//map[<=30:[Bob Tom] >30:[Alice Chris]]

	assert.Equal(t, []string{"Bob", "Tom"}, ageGroupedSortedNames["<=30"])
	assert.Equal(t, []string{"Alice", "Chris"}, ageGroupedSortedNames[">30"])
}
