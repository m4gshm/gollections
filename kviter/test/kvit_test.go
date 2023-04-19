package test

import (
	"errors"
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/kviter"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Collect_Group(t *testing.T) {
	groups := kviter.Group(kviter.OfPairs(K.V(1, "1"), K.V(2, "2"), K.V(2, "22")))

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []string{"1"}, groups[1])
	assert.Equal(t, []string{"2", "22"}, groups[2])
}

func Test_Collect_Map(t *testing.T) {
	groups := kviter.ToMap(kviter.FromPairs(iter.Of(K.V(1, "1"), K.V(2, "2"), K.V(2, "22"))))

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, "1", groups[1])
	assert.Equal(t, "2", groups[2])
}

type rows[T any] struct {
	in     []T
	cursor int
}

func (r *rows[T]) hasNext() bool {
	return r.cursor < len(r.in)
}

func (r *rows[T]) next() (T, error) {
	e := r.in[r.cursor]
	r.cursor++
	if r.cursor > 3 {
		var no T
		return no, errors.New("next error")
	}
	return e, nil
}

func Test_OfLoop(t *testing.T) {

	stream := &rows[int]{slice.Of(1, 2, 3), 0}

	evens := func(r *rows[int]) (bool, int, error) {
		next, err := r.next()
		if err != nil {
			return false, 0, err
		}
		return next%2 == 0, next, nil
	}

	iterator := kviter.OfLoop(stream, (*rows[int]).hasNext, evens)

	m := kviter.ToMap[bool, int](iterator)

	assert.Equal(t, 2, m[true])
	assert.Equal(t, 1, m[false])
	assert.Nil(t, iterator.Error())

	streamWithError := &rows[int]{slice.Of(1, 2, 3, 4), 0}
	iterWithError := kviter.OfLoop(streamWithError, (*rows[int]).hasNext, evens)
	m2 := kviter.ToMap[bool, int](iterWithError)

	assert.Equal(t, 2, m2[true])
	assert.Equal(t, "next error", iterWithError.Error().Error())
}
