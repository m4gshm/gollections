package iter

import (
	"github.com/m4gshm/gollections/c"
)

// FitKV is the KVIterator wrapper that provides filtering of key/value elements by a Predicate.
type FitKV[K, V any] struct {
	next func() (K, V, bool)
	by   func(K, V) bool
}

var (
	_ c.KVIterator[any, any] = (*FitKV[any, any])(nil)
	_ c.KVIterator[any, any] = FitKV[any, any]{}
)

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f FitKV[K, V]) Next() (key K, value V, ok bool) {
	if !(f.next == nil || f.by == nil) {
		key, value, ok = nextFilteredKV(f.next, f.by)
	}
	return key, value, ok
}

func nextFilteredKV[K any, V any](next func() (K, V, bool), filter func(K, V) bool) (key K, val V, filtered bool) {
	for key, val, ok := next(); ok; key, val, ok = next() {
		if filter(key, val) {
			return key, val, true
		}
	}
	return key, val, false
}
