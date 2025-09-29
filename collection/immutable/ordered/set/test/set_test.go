package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Set_From(t *testing.T) {
	set := set.From(loop.Of(1, 1, 2, 2, 3, 4, 3, 2, 1))
	assert.Equal(t, slice.Of(1, 2, 3, 4), set.Slice())
}

func Test_Set_FromSeq(t *testing.T) {
	set := set.FromSeq(seq.Of(1, 1, 2, 2, 3, 4, 3, 2, 1))
	assert.Equal(t, slice.Of(1, 2, 3, 4), set.Slice())
}

func Test_Set_Iterate(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := set.Slice()

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 4, 3)
	assert.Equal(t, expected, values)

	loopService := seq.Slice(set.All)
	assert.Equal(t, expected, loopService)

	out := make([]int, 0)
	for v := range set.All {
		out = append(out, v)
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	set.ForEach(func(v int) { out = append(out, v) })

	assert.Equal(t, expected, out)
}

func Test_Set_Contains(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(4))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(0))
	assert.False(t, set.Contains(-1))
}

func Test_Set_FilterReduce(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1).Reduce(op.Sum[int])
	assert.Equal(t, 1+2+3+4, s)
}

func Test_Set_FilterMapReduce(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, s)
}

func Test_Set_Group_By_Walker(t *testing.T) {
	groups := group.Of(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_Set_Group_By_Iterator(t *testing.T) {
	groups := seq.Group(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7).All, func(e int) bool { return e%2 == 0 }, as.Is[int])

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_Set_Sort(t *testing.T) {
	var (
		elements = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		sorted   = elements.Sort(op.Compare)
	)
	assert.Equal(t, set.Of(-2, 0, 1, 3, 5, 6, 8), sorted)
	assert.NotSame(t, &elements, &sorted)
}
func Test_Set_SortStructByField(t *testing.T) {
	var (
		anonymous = &user{"Anonymous", 0}
		cherlie   = &user{"Cherlie", 25}
		alise     = &user{"Alise", 20}
		bob       = &user{"Bob", 19}

		elements     = set.Of(anonymous, cherlie, alise, bob)
		sortedByName = set.Sort(elements, (*user).Name)
		sortedByAge  = set.Sort(elements, (*user).Age)
	)
	assert.Equal(t, set.Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, set.Of(anonymous, bob, alise, cherlie), sortedByAge)
}

func Test_Set_Convert(t *testing.T) {
	var (
		ints     = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		strings  = seq.Slice(seq.Filter(set.Convert(ints, strconv.Itoa), func(s string) bool { return len(s) == 1 }))
		strings2 = set.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }).Slice()
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), strings)
	assert.Equal(t, strings, strings2)
}

func Test_Set_Flatt(t *testing.T) {
	var (
		ints        = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		fints       = set.Flat(ints, func(i int) []int { return slice.Of(i) })
		stringsPipe = seq.Filter(seq.Convert(fints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }), func(s string) bool { return len(s) == 1 })
	)
	assert.Equal(t, slice.Of("3", "1", "5", "6", "8", "0"), stringsPipe.Slice())
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }
