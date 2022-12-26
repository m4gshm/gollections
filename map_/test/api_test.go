package test

import (
	"testing"

	"github.com/m4gshm/gollections/kvit"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/group"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_StringRepresentation(t *testing.T) {
	order := slice.Of(4, 3, 2, 1)
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}
	actual := map_.ToStringOrdered(order, elements)

	expected := "[4:4 3:3 2:2 1:1]"
	assert.Equal(t, expected, actual)
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
	return e, nil
}

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := map_.OfLoop(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	})

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_OfLoopResolv(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3, 4), 0}
	result, _ := map_.OfLoopResolv(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	}, kvit.LastVal[bool, int])

	assert.Equal(t, 4, result[true])
	assert.Equal(t, 3, result[false])
}

func Test_Generate(t *testing.T) {
	counter := 0
	result, _ := map_.Generate(func() (bool, int, bool, error) {
		counter++
		return counter%2 == 0, counter, counter < 4, nil
	})

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_GenerateResolv(t *testing.T) {
	counter := 0
	result, _ := map_.GenerateResolv(func() (bool, int, bool, error) {
		counter++
		return counter%2 == 0, counter, counter < 5, nil
	}, func(exists bool, k bool, old, new int) int {
		return op.IfElse(exists, op.IfElse(k, new, old), new)
	})

	assert.Equal(t, 4, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_GroupOfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := group.OfLoop(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	})

	assert.Equal(t, slice.Of(2), result[true])
	assert.Equal(t, slice.Of(1, 3), result[false])
}
