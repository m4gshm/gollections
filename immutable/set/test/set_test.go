package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Set_From(t *testing.T) {
	set := set.From(iter.Of(1, 1, 2, 2, 3, 4, 3, 2, 1).Next)
	assert.Equal(t, slice.Of(1, 2, 3, 4), sort.Of(set.Slice()))
}

func Test_Set_Iterate(t *testing.T) {
	set := set.Of(1, 1, 2, 4, 3, 1)
	values := sort.Of(set.Slice())

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 3, 4)
	assert.Equal(t, expected, values)

	loopSlice := sort.Of(loop.Slice(ptr.Of(set.Head()).Next))
	assert.Equal(t, expected, loopSlice)

	out := make(map[int]int, 0)
	for it, v, ok := set.First(); ok; v, ok = it.Next() {
		out[v] = v
	}

	assert.Equal(t, len(expected), len(out))
	for k := range out {
		assert.True(t, set.Contains(k))
	}

	out = make(map[int]int, 0)
	set.ForEach(func(n int) { out[n] = n })

	assert.Equal(t, len(expected), len(out))
	for k := range out {
		assert.True(t, set.Contains(k))
	}
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

func Test_Set_FilterMapReduce(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, s)
}

func Test_Set_Group_By_Walker(t *testing.T) {
	groups := group.Of(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	fg := sort.Of(groups[false])
	tg := sort.Of(groups[true])
	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, fg)
	assert.Equal(t, []int{0, 2, 4, 6}, tg)
}

func Test_Set_Group_By_Iterator(t *testing.T) {
	groups := loop.Group(set.Of(0, 1, 1, 2, 4, 3, 1, 6, 7).Iter().Next, func(e int) bool { return e%2 == 0 }, as.Is[int])

	assert.Equal(t, len(groups), 2)
	fg := sort.Of(groups[false])

	tg := sort.Of(groups[true])

	assert.Equal(t, []int{1, 3, 7}, fg)
	assert.Equal(t, []int{0, 2, 4, 6}, tg)
}

func Test_Set_Sort(t *testing.T) {
	var (
		elements = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, oset.Of(-2, 0, 1, 3, 5, 6, 8), sorted)
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
	assert.Equal(t, oset.Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, oset.Of(anonymous, bob, alise, cherlie), sortedByAge)
}

func Test_Set_Convert(t *testing.T) {
	var (
		ints     = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		strings  = sort.Of(loop.Slice(iter.Filter(set.Convert(ints, strconv.Itoa), func(s string) bool { return len(s) == 1 }).Next))
		strings2 = sort.Of(set.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 }).Slice())
	)
	assert.Equal(t, slice.Of("0", "1", "3", "5", "6", "8"), strings)
	assert.Equal(t, strings, strings2)
}

func Test_Set_Flatt(t *testing.T) {

	var (
		ints        = set.Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		fints       = set.Flatt(ints, func(i int) []int { return slice.Of(i) })
		convFilt    = collection.Convert(fints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 })
		stringsPipe = collection.Filter(convFilt, func(s string) bool { return len(s) == 1 })
	)
	assert.Equal(t, slice.Of("0", "1", "3", "5", "6", "8"), sort.Of(stringsPipe.Slice()))
}

func Test_Set_DoubleConvert(t *testing.T) {
	var (
		ints               = set.Of(3, 1, 5, 6, 8, 0, -2)
		stringsPipe        = set.Convert(ints, strconv.Itoa).Filter(func(s string) bool { return len(s) == 1 })
		prefixedStrinsPipe = collection.Convert(stringsPipe, func(s string) string { return "_" + s })
	)
	assert.Equal(t, slice.Of("_0", "_1", "_3", "_5", "_6", "_8"), sort.Of(prefixedStrinsPipe.Slice()))

	//second call do nothing
	var no []string
	assert.Equal(t, no, stringsPipe.Slice())
}

func Test_Set_Zero(t *testing.T) {
	var set immutable.Set[int]

	assert.False(t, set.Contains(1))

	set.IsEmpty()
	set.Len()

	set.For(nil)
	set.ForEach(nil)

	set.Slice()

	set.Convert(nil)
	set.Filter(nil)

	head := set.Head()
	_, ok := head.Next()
	assert.False(t, ok)

	_, _, ok = set.First()
	assert.False(t, ok)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }
