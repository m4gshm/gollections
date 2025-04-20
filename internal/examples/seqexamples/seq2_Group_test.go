package seqexamples

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/expr/use"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
)

func Test_Group(t *testing.T) {

	var users iter.Seq[User] = seq.Of(users...)
	var groups iter.Seq2[string, User] = seq.ToSeq2(users, func(u User) (string, User) {
		return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30"), u
	})
	var ageGroups map[string][]User = seq2.Group(groups)

	//map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]

	assert.Equal(t, slice.Of("Alice", "Chris"), sort.Asc(slice.Convert(ageGroups[">30"], User.Name)))
	assert.Equal(t, slice.Of("Bob"), slice.Convert(ageGroups["<=30"], User.Name))
	assert.Equal(t, slice.Of("Tom"), slice.Convert(ageGroups["<=20"], User.Name))
}
