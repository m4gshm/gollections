// Package iter provides ordered map iterator implementations
package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// NewOrdered is the OrderedKV constructor
func NewOrdered[K comparable, V any](uniques map[K]V, elements slice.Iter[K]) OrderedMapIter[K, V] {
	return OrderedMapIter[K, V]{elements: elements, uniques: uniques}
}

// OrderedMapIter is the ordered key/value pairs Iterator implementation
type OrderedMapIter[K comparable, V any] struct {
	elements slice.Iter[K]
	uniques  map[K]V
}

var _ kv.KVIterator[string, any] = (*OrderedMapIter[string, any])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *OrderedMapIter[K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i *OrderedMapIter[K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(i.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *OrderedMapIter[K, V]) Next() (key K, val V, ok bool) {
	if i != nil {
		if key, ok = i.elements.Next(); ok {
			val = i.uniques[key]
		}
	}
	return key, val, ok
}

// Cap returns the iterator capacity
func (i *OrderedMapIter[K, V]) Cap() int {
	return i.elements.Cap()
}

// NewValIter is default ValIter constructor
func NewValIter[K comparable, V any](elements []K, uniques map[K]V) ValIter[K, V] {
	return ValIter[K, V]{elements: elements, uniques: uniques, current: slice.IterNoStarted}
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

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *ValIter[K, V]) For(walker func(element V) error) error {
	return loop.For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *ValIter[K, V]) ForEach(walker func(element V)) {
	loop.ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *ValIter[K, V]) Next() (val V, ok bool) {
	if i != nil && slice.HasNext(i.elements, i.current) {
		i.current++
		return i.uniques[slice.Get(i.elements, i.current)], true
	}
	return val, false
}

// Cap returns the iterator capacity
func (i *ValIter[K, V]) Cap() int {
	if i == nil {
		return 0
	}
	return len(i.elements)
}
