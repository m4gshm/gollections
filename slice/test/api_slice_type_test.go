package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	_less "github.com/m4gshm/gollections/break/predicate/less"
	_more "github.com/m4gshm/gollections/break/predicate/more"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
)

func Test_Slice_Clone(t *testing.T) {
	type entity struct{ val string }
	var (
		first  = entity{"first"}
		second = entity{"second"}
		third  = entity{"third"}

		entities = slice.Wrap(&first, &second, &third)
		c        = entities.Clone()
	)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Same(t, entities[i], c[i])
	}
}

func Test_Slice_DeepClone(t *testing.T) {
	type entity struct{ val string }
	var (
		first  = entity{"first"}
		second = entity{"second"}
		third  = entity{"third"}

		entities = slice.Wrap(&first, &second, &third)
		c        = entities.DeepClone(clone.Ptr[entity])
	)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Equal(t, entities[i], c[i])
		assert.NotSame(t, entities[i], c[i])
	}
}

func Test_Slice_ReduceSum(t *testing.T) {
	s := slice.Wrap(1, 3, 5, 7, 9, 11)
	r := s.Reduce(op.Sum)
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Slice_First(t *testing.T) {
	s := slice.Wrap(1, 3, 5, 7, 9, 11)
	r, ok := s.First(func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := s.First(func(i int) bool { return i > 12 })
	assert.False(t, nook)
}

func Test_Slice_Firstt(t *testing.T) {
	s := slice.Wrap(1, 3, 5, 7, 9, 11)
	r, ok, _ := s.Firstt(_more.Than(5))
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook, _ := s.Firstt(_more.Than(12))
	assert.False(t, nook)

	_, _, err := s.Firstt(func(i int) (bool, error) { return true, errors.New("abort") })
	assert.Error(t, err)
}

func Test_Slice_Last(t *testing.T) {
	s := slice.Wrap(1, 3, 5, 7, 9, 11)
	r, ok := s.Last(func(i int) bool { return i < 9 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := s.Last(func(i int) bool { return i < 1 })
	assert.False(t, nook)
}

func Test_Slice_Lastt(t *testing.T) {
	s := slice.Wrap(1, 3, 5, 7, 9, 11)
	r, ok, _ := s.Lastt(_less.Than(9))
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook, _ := s.Lastt(_less.Than(1))
	assert.False(t, nook)

	_, _, err := s.Lastt(func(i int) (bool, error) { return true, errors.New("abort") })
	assert.Error(t, err)
}

func Test_Slice_Filter(t *testing.T) {
	s := slice.Wrap(1, 3, 4, 5, 7, 8, 9, 11)
	r := s.Filter(even)
	assert.Equal(t, slice.Wrap(4, 8), r)
}

func Test_Slice_Filt(t *testing.T) {
	var s slice.Slice[int] = []int{1, 3, 4, 5, 7, 8, 9, 11}
	r, err := s.Filt(func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) })
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4, 8), r.Unwrap())
}

func Test_Slice_Filt2(t *testing.T) {
	s := slice.Wrap(1, 3, 4, 5, 7, 8, 9, 11)
	r, err := s.Filt(func(i int) (bool, error) {
		ok := i <= 7
		return ok && even(i), op.IfElse(ok, nil, errors.New("abort"))

	})
	assert.Error(t, err)
	assert.Equal(t, slice.Wrap(4), r)
}

func Test_Slice_MatchAny(t *testing.T) {
	s := slice.Wrap(1, 2, 3, 4)

	ok := s.HasAny(eq.To(4))
	assert.True(t, ok)

	noOk := s.HasAny(more.Than(5))
	assert.False(t, noOk)
}

func Test_Slice_Empty(t *testing.T) {
	assert.False(t, slice.Wrap(1).Empty())
	assert.True(t, slice.Wrap[int]().Empty())

	var s slice.Slice[int]
	assert.True(t, s.Empty())
}
