package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

//NewValIter is default ValIter constructor.
func NewValIter[K comparable, V any](elements []K, uniques map[K]V) *ValIter[K, V] {
	return &ValIter[K, V]{elements: elements, uniques: uniques}
}

//ValIter is the Iteratoc over Map values.
type ValIter[K comparable, V any] struct {
	elements []K
	uniques  map[K]V
	current  int
}

var _ c.Iterator[any] = (*ValIter[int, any])(nil)

func (s *ValIter[K, V]) HasNext() bool {
	return it.HasNext(s.elements, s.current)
}

func (s *ValIter[K, V]) Next() V {
	if s.HasNext() {
		s.current++
		return s.uniques[it.Get(s.elements, s.current)]
	}
	var no V
	return no
}
