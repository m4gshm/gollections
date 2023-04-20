package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// Fit is the Iterator wrapper that provides filtering of elements by a Predicate.
type Fit[T, IT any] struct {
	iterator IT
	next     func() (T, bool)
	by       func(T) bool
}

var (
	_ c.Iterator[any] = (*Fit[any, any])(nil)
	_ c.Iterator[any] = Fit[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s Fit[T, IT]) Next() (T, bool) {
	return nextFiltered(s.next, s.by)
}

// FitKV is the KVIterator wrapper that provides filtering of key/value elements by a Predicate.
type FitKV[K, V any, IT c.KVIterator[K, V]] struct {
	iterator IT
	by       func(K, V) bool
}

var (
	_ c.KVIterator[any, any] = (*FitKV[any, any, c.KVIterator[any, any]])(nil)
	_ c.KVIterator[any, any] = FitKV[any, any, c.KVIterator[any, any]]{}
)

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s FitKV[K, V, IT]) Next() (K, V, bool) {
	return nextFilteredKV(s.iterator.Next, s.by)
}

func nextFiltered[T any](next func() (T, bool), filter func(T) bool) (v T, ok bool) {
	return loop.First(next, filter)
}

func nextFilteredKV[K any, V any](next func() (K, V, bool), filter func(K, V) bool) (key K, val V, filtered bool) {
	for key, val, ok := next(); ok; key, val, ok = next() {
		if filter(key, val) {
			return key, val, true
		}
	}
	return key, val, false
}
