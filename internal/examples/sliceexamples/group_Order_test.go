package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/expr/use"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/m4gshm/gollections/slice/sort"
)

func Test_Slice_GroupOrder(t *testing.T) {

	var order, ageGroups = group.Order(users, func(u User) string {
		return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30")
	}, as.Is)

	//order     [<=30 >30 <=20]
	//ageGroups map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]

	assert.Equal(t, slice.Of("Alice", "Chris"), sort.Asc(slice.Convert(ageGroups[">30"], User.Name)))
	assert.Equal(t, slice.Of("Bob"), slice.Convert(ageGroups["<=30"], User.Name))
	assert.Equal(t, slice.Of("Tom"), slice.Convert(ageGroups["<=20"], User.Name))
	assert.Equal(t, slice.Of("<=30", ">30", "<=20"), order)
}
