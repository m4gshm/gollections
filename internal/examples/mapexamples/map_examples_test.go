package mapexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/clone"
	"github.com/m4gshm/gollections/map_/group"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone/sort"
)

type entity struct{ val string }

var (
	first  = entity{"1_first"}
	second = entity{"2_second"}
	third  = entity{"3_third"}

	entities = map[int]*entity{1: &first, 2: &second, 3: &third}
)

func Test_DeepClone(t *testing.T) {
	c := clone.Deep(entities, func(e *entity) *entity { return ptr.Of(*e) })

	assert.Equal(t, entities, c)
	assert.NotSame(t, &entities, &c)

	for i := range entities {
		assert.Equal(t, entities[i], c[i])
		assert.NotSame(t, entities[i], c[i])
	}
}

func Test_ValuesConverted(t *testing.T) {
	var values []string = map_.ValuesConverted(entities, func(e *entity) string { return e.val })
	assert.Equal(t, slice.Of("1_first", "2_second", "3_third"), sort.Asc(values))
}

type rows[T any] struct {
	in     []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.in) }
func (r *rows[T]) next() (T, error) { e := r.in[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := map_.OfLoop(
		stream,
		(*rows[int]).hasNext,
		func(r *rows[int]) (bool, int, error) {
			n, err := r.next()
			return n%2 == 0, n, err
		},
	)

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
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

func Test_GroupOfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := group.OfLoop(
		stream,
		(*rows[int]).hasNext,
		func(r *rows[int]) (bool, int, error) {
			n, err := r.next()
			return n%2 == 0, n, err
		},
	)

	assert.Equal(t, slice.Of(2), result[true])
	assert.Equal(t, slice.Of(1, 3), result[false])
}
