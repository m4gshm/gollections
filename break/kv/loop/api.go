// Package loop provides helpers for loop operation over key/value pairs and iterator implementations
package loop

import (
	"github.com/m4gshm/gollections/map_/resolv"
)

// From wrap the next loop to a breakable loop
func From[K, V any](next func() (K, V, bool)) func() (K, V, bool, error) {
	return func() (K, V, bool, error) {
		k, v, ok := next()
		return k, v, ok, nil
	}
}

// To transforms a breakable loop to a simple loop.
// The errConsumer is a function that is called when an error occurs.
func To[K, V any](next func() (K, V, bool, error), errConsumer func(error)) func() (K, V, bool) {
	return func() (K, V, bool) {
		k, v, ok, err := next()
		if err != nil {
			errConsumer(err)
			return k, v, false
		}
		return k, v, ok
	}
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](next func() (K, V, bool, error)) (map[K][]V, error) {
	return ToMapResolv(next, resolv.Append[K, V])
}

// Reduce reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function
func Reduce[K, V any](next func() (K, V, bool, error), merge func(K, V, K, V) (K, V, error)) (rk K, rv V, err error) {
	k, v, ok, err := next()
	if err != nil || !ok {
		return rk, rv, err
	}
	rk, rv = k, v
	for {
		if k, v, ok, err := next(); err != nil || !ok {
			return rk, rv, err
		} else if rk, rv, err = merge(rk, rv, k, v); err != nil {
			return rk, rv, err
		}
	}
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAny[K, V any](next func() (K, V, bool, error), predicate func(K, V) (bool, error)) (bool, error) {
	for {
		if k, v, ok, err := next(); err != nil || !ok {
			return false, err
		} else if ok, err := predicate(k, v); err != nil {
			return false, err
		} else if ok {
			return true, nil
		}
	}
}

// Convert creates an iterator that applies a transformer to iterable key\values.
func Convert[K, V any, KOUT, VOUT any](next func() (K, V, bool, error), converter func(K, V) (KOUT, VOUT)) ConvertIter[K, V, KOUT, VOUT] {
	return ConvertIter[K, V, KOUT, VOUT]{next: next, converter: func(k K, v V) (KOUT, VOUT, error) { ko, vo := converter(k, v); return ko, vo, nil }}
}

// Conv creates an iterator that applies a transformer to iterable key\values.
func Conv[K, V any, KOUT, VOUT any](next func() (K, V, bool, error), converter func(K, V) (KOUT, VOUT, error)) ConvertIter[K, V, KOUT, VOUT] {
	return ConvertIter[K, V, KOUT, VOUT]{next: next, converter: converter}
}

// Filter creates an iterator that checks elements by a filter and returns successful ones
func Filter[K, V any](next func() (K, V, bool, error), filter func(K, V) bool) FiltKV[K, V] {
	return FiltKV[K, V]{next: next, filter: func(k K, v V) (bool, error) { return filter(k, v), nil }}
}

// Filt creates an iterator that checks elements by a filter and returns successful ones
func Filt[K, V any](next func() (K, V, bool, error), filter func(K, V) (bool, error)) FiltKV[K, V] {
	return FiltKV[K, V]{next: next, filter: filter}
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[K comparable, V, VR any](next func() (K, V, bool, error), resolver func(bool, K, VR, V) VR) (map[K]VR, error) {
	m := map[K]VR{}
	for {
		k, v, ok, err := next()
		if err != nil || !ok {
			return m, err
		}
		exists, ok := m[k]
		m[k] = resolver(ok, k, exists, v)
	}
}

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[K comparable, V any](next func() (K, V, bool, error)) (map[K]V, error) {
	return ToMapResolv(next, resolv.FirstVal[K, V])
}

// ToSlice collects key\value elements to a slice by iterating over the elements
func ToSlice[K, V, T any](next func() (K, V, bool, error), converter func(K, V) T) ([]T, error) {
	s := []T{}
	for {
		key, val, ok, err := next()
		if err != nil || !ok {
			return s, err
		}
		s = append(s, converter(key, val))
	}
}

// New is the mai breakable key/value loop constructor
func New[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) func() (K, V, bool, error) {
	return func() (k K, v V, ok bool, err error) {
		if ok := hasNext(source); !ok {
			return k, v, false, nil
		} else if k, v, err = getNext(source); err != nil {
			return k, v, false, err
		} else {
			return k, v, true, nil
		}
	}
}
