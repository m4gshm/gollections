package vector

import (
	"errors"
	"sync"
	"testing"

	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/sum"
	"github.com/m4gshm/gollections/walk/group"

	"github.com/stretchr/testify/assert"
)

func Test_Vector_Iterate(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	values := vec.Collect()

	assert.Equal(t, 6, len(values))

	expected := slice.Of(1, 1, 2, 4, 3, 1)
	assert.Equal(t, expected, values)

	iterSlice := it.Slice(vec.Begin())
	assert.Equal(t, expected, iterSlice)

	out := make([]int, 0)
	for it := vec.Begin(); it.HasNext(); {
		n, _ := it.Get()
		out = append(out, n)
	}
	assert.Equal(t, expected, out)

	out = make([]int, 0)
	vec.ForEach(func(v int) { out = append(out, v) })
}

func Test_Vector_Add(t *testing.T) {
	vec := vector.New[int](0)
	added, _ := vec.Add(1, 1, 2, 4, 3, 1)
	assert.Equal(t, added, true)
	added, _ = vec.Add(1)
	assert.Equal(t, added, true)
	assert.Equal(t, slice.Of(1, 1, 2, 4, 3, 1, 1), vec.Collect())
}

func Test_Vector_DeleteOne(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	r, _ := vec.Delete(3)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("1", "1", "2", "3", "1"), vec.Collect())
	r, _ = vec.Delete(5)
	assert.Equal(t, r, false)
}

func Test_Vector_DeleteMany(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r, _ := vec.Delete(3, 0, 5)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("1", "2", "4", "6"), vec.Collect())
	r, _ = vec.Delete(5, 4)
	assert.Equal(t, r, false)
}

func Test_Vector_DeleteManyFromTail(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r, _ := vec.Delete(4, 5, 6)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("0", "1", "2", "3"), vec.Collect())
}

func Test_Vector_DeleteManyFromHead(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r, _ := vec.Delete(0, 1, 2)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("3", "4", "5", "6"), vec.Collect())
}

func Test_Vector_DeleteManyFromMiddle(t *testing.T) {
	vec := vector.Of("0", "1", "2", "3", "4", "5", "6")
	r, _ := vec.Delete(4, 3)
	assert.Equal(t, r, true)
	assert.Equal(t, slice.Of("0", "1", "2", "5", "6"), vec.Collect())
}

func Test_Vector_Set(t *testing.T) {
	vec := vector.Of("1", "1", "2", "4", "3", "1")
	added, _ := vec.Set(10, "11")
	assert.Equal(t, added, true)
	assert.Equal(t, slice.Of("1", "1", "2", "4", "3", "1", "", "", "", "", "11"), vec.Collect())
}

func Test_Vector_DeleteByIterator(t *testing.T) {
	vec := vector.Of(1, 1, 2, 4, 3, 1)
	iter := vec.BeginEdit()

	i := 0
	for iter.HasNext() {
		i++
		_, _ = iter.Delete()
	}

	assert.Equal(t, 6, i)
	assert.Equal(t, 0, len(vec.Collect()))
}

func Test_Vector_FilterMapReduce(t *testing.T) {
	s := vector.Of(1, 1, 2, 4, 3, 4).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
	assert.Equal(t, 20, s)

	s = it.Pipe(vector.Of(1, 1, 2, 4, 3, 1, 4).Begin()).Filter(func(i int) bool { return i%2 == 0 }).Map(func(i int) int { return i * 2 }).Reduce(sum.Of[int])
	assert.Equal(t, 20, s)
}

func Test_Vector_Group(t *testing.T) {
	groups := group.Of(vector.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 })

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 1, 3, 1, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_Vector_Concurrent_Update(t *testing.T) {
	vec := vector.Empty[int64]()

	wg := sync.WaitGroup{}

	wg.Add(1)
	var add error
	go func() {
		defer wg.Done()
		i := int64(0)
		for {
			if _, err := vec.Add(i); err != nil {
				if errors.Is(err, mutable.BadRW) {
					add = err
					return
				}
			}
			i++
		}
	}()

	wg.Add(1)
	var delete error
	go func() {
		defer wg.Done()
		for {
			for iter := vec.BeginEdit(); iter.HasNext(); {
				if _, err := iter.Delete(); err != nil {
					if errors.Is(err, mutable.BadRW) {
						delete = err
						return
					}

				}
			}
		}
	}()

	wg.Wait()
	if add == nil || delete == nil {
		t.Fatal("no errors")
	}
}
