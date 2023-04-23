package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/iter"
)

// NewValIter is default ValIter constructor
func NewValIter[K comparable, V any](elements []K, uniques map[K]V) ValIter[K, V] {
	return ValIter[K, V]{elements: elements, uniques: uniques, current: iter.NoStarted}
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

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (v *ValIter[K, V]) Next() (val V, ok bool) {
	if v != nil && iter.HasNext(v.elements, v.current) {
		v.current++
		return v.uniques[iter.Get(v.elements, v.current)], true
	}
	return val, false
}

// Cap returns the iterator capacity
func (v *ValIter[K, V]) Cap() int {
	if v == nil {
		return 0
	}
	return len(v.elements)
}
