package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// NewValIter is default ValIter constructor
func NewValIter[K comparable, V any](elements []K, uniques map[K]V) *ValIter[K, V] {
	return &ValIter[K, V]{elements: elements, uniques: uniques, current: iter.NoStarted}
}

// ValIter is the Iteratoc over Map values
type ValIter[K comparable, V any] struct {
	elements []K
	uniques  map[K]V
	current  int
}

var (
	_ c.Iterator[any] = (*ValIter[int, any])(nil)
	_ c.Sized         = (*ValIter[int, any])(nil)
)

func (s *ValIter[K, V]) Next() (V, bool) {
	if iter.HasNext(s.elements, s.current) {
		s.current++
		return s.uniques[iter.Get(s.elements, s.current)], true
	}
	var no V
	return no, false
}

func (s *ValIter[K, V]) Cap() int {
	return len(s.elements)
}
