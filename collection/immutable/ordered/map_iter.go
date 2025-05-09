// Package ordered provides ordered map iterator implementations
package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv/collection"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// NewMapIter is the Iter constructor
func NewMapIter[K comparable, V any](uniques map[K]V, elements slice.Iter[K]) MapIter[K, V] {
	return MapIter[K, V]{elements: elements, uniques: uniques}
}

// MapIter is the ordered key/value pairs Iterator implementation
type MapIter[K comparable, V any] struct {
	elements slice.Iter[K]
	uniques  map[K]V
}

var _ collection.Iterator[string, any] = (*MapIter[string, any])(nil)

// All is used to iterate through the iterator using `for ... range`.
func (i *MapIter[K, V]) All(consumer func(key K, value V) bool) {
	kvloop.All(i.Next, consumer)
}

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning Break
func (i *MapIter[K, V]) Track(traker func(key K, value V) error) error {
	return kvloop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i *MapIter[K, V]) TrackEach(traker func(key K, value V)) {
	kvloop.TrackEach(i.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *MapIter[K, V]) Next() (key K, val V, ok bool) {
	if i != nil {
		if key, ok = i.elements.Next(); ok {
			val = i.uniques[key]
		}
	}
	return key, val, ok
}

// Size returns the iterator capacity
func (i *MapIter[K, V]) Size() int {
	return i.elements.Size()
}

// NewValIter is default ValIter constructor
func NewValIter[K comparable, V any](elements []K, uniques map[K]V) *ValIter[K, V] {
	return &ValIter[K, V]{elements: elements, uniques: uniques, current: slice.IterNoStarted}
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

// All is used to iterate through the iterator using `for ... range`.
func (i *ValIter[K, V]) All(consumer func(element V) bool) {
	loop.All(i.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning Break
func (i *ValIter[K, V]) For(consumer func(element V) error) error {
	return loop.For(i.Next, consumer)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *ValIter[K, V]) ForEach(consumer func(element V)) {
	loop.ForEach(i.Next, consumer)
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

// Size returns the iterator capacity
func (i *ValIter[K, V]) Size() int {
	if i == nil {
		return 0
	}
	return len(i.elements)
}
