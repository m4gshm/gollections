package oset

import (
	"testing"

	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/sum"

	"github.com/m4gshm/gollections/walk/group"
	"github.com/stretchr/testify/assert"
)

func Test_Set_Iterate(t *testing.T) {
	set := Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 4, 3)
	assert.Equal(t, expected, values)

	iterSlice := it.Slice(set.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for it := set.Begin(); it.HasNext(); {
		n := it.Next()
		out = append(out, n)
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	set.ForEach(func(v int) { out = append(out, v) })

	assert.Equal(t, expected, out)
}

func Test_Set_Contains(t *testing.T) {
	set := Of(1, 1, 2, 4, 3, 1)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(4))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(0))
	assert.False(t, set.Contains(-1))
}

func Test_Set_FilterMapReduce(t *testing.T) {
	s := Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
	assert.Equal(t, 12, s)

	s = it.Pipe(Of(1, 1, 2, 4, 3, 1).Begin()).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
	assert.Equal(t, 12, s)
}

func Test_Set_Group_By_Walker(t *testing.T) {
	groups := group.Of(Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_Set_Group_By_Iterator(t *testing.T) {
	groups := it.Group(Of(0, 1, 1, 2, 4, 3, 1, 6, 7).Begin(), func(e int) bool { return e%2 == 0 }).Collect()

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_Set_Sort(t *testing.T) {
	var (
		elements = Of(3, 3, 1, 1, 1, 5, 6, 8, 8, 0, -2, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, Of(-2, 0, 1, 3, 5, 6, 8), sorted)
}
func Test_Set_SortStructByField(t *testing.T) {
	var (
		anonymous = &user{"Anonymous", 0}
		cherlie   = &user{"Cherlie", 25}
		alise     = &user{"Alise", 20}
		bob       = &user{"Bob", 19}

		elements     = Of(anonymous, cherlie, alise, bob)
		sortedByName = Sort(elements, (*user).Name)
		sortedByAge  = Sort(elements, (*user).Age)
	)
	assert.Equal(t, Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, Of(anonymous, bob, alise, cherlie), sortedByAge)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }
