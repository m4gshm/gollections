package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/mutable/oset"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Set_Iterate(t *testing.T) {
	set := oset.Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	assert.Equal(t, 4, len(values))

	expected := slice.Of(1, 2, 4, 3)
	assert.Equal(t, expected, values)

	iterSlice := it.ToSlice(set.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	it := set.Begin()
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		out = append(out, v)
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	set.ForEach(func(v int) { out = append(out, v) })
}

func Test_Set_AddVerify(t *testing.T) {
	set := oset.New[int](0)
	added := set.AddNew(1, 2, 4, 3)
	assert.Equal(t, added, true)
	added = set.AddNewOne(1)
	assert.Equal(t, added, false)

	values := set.Collect()

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_Set_Delete(t *testing.T) {
	set := oset.Of(1, 1, 2, 4, 3, 1)
	values := set.Collect()

	for _, v := range values {
		set.Delete(v)
	}

	assert.Equal(t, 0, len(set.Collect()))
}

func Test_Set_DeleteByIterator(t *testing.T) {
	set := oset.Of(1, 1, 2, 4, 3, 1)
	iter := set.BeginEdit()

	i := 0
	for _, ok := iter.Next(); ok; _, ok = iter.Next() {
		i++
		iter.Delete()
	}

	assert.Equal(t, 4, i)
	assert.Equal(t, 0, len(set.Collect()))
}

func Test_Set_FilterMapReduce(t *testing.T) {
	s := oset.Of(1, 1, 2, 4, 3, 1).Filter(func(i int) bool { return i%2 == 0 }).Convert(func(i int) int { return i * 2 }).Reduce(op.Sum[int])
	assert.Equal(t, 12, s)
}

func Test_Set_Group(t *testing.T) {
	groups := group.Of(oset.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 3, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}
