// Package loop provides helpers for loop operation over key/value pairs and iterator implementations
package loop

import (
	"github.com/m4gshm/gollections/break/kv/loop"
	"github.com/m4gshm/gollections/map_/resolv"
)

// Looper provides an iterable loop function
type Looper[K, V any, I interface{ Next() (K, V, bool) }] interface {
	Loop() I
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](next func() (K, V, bool)) map[K][]V {
	return ToMapResolv(next, resolv.Append[K, V])
}

// Reduce reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function
func Reduce[K, V any](next func() (K, V, bool), merge func(K, K, V, V) (K, V)) (rk K, rv V) {
	if k, v, ok := next(); ok {
		rk, rv = k, v
	} else {
		return rk, rv
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		rk, rv = merge(rk, k, rv, v)
	}
	return rk, rv
}

// Reducee reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function
func Reducee[K, V any](next func() (K, V, bool), merge func(K, K, V, V) (K, V, error)) (rk K, rv V, err error) {
	k, v, ok := next()
	if !ok {
		return rk, rv, nil
	}
	rk, rv = k, v
	for {
		if k, v, ok := next(); !ok {
			return rk, rv, nil
		} else if rk, rv, err = merge(rk, k, rv, v); err != nil {
			return rk, rv, err
		}
	}
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAny[K, V any](next func() (K, V, bool), predicate func(K, V) bool) bool {
	_, _, ok := First(next, predicate)
	return ok
}

// HasAnyy finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAnyy[K, V any](next func() (K, V, bool), predicate func(K, V) (bool, error)) (bool, error) {
	_, _, ok, err := Firstt(next, predicate)
	return ok, err
}

// First returns the first key/value pair that satisfies the condition of the 'predicate' function
func First[K, V any](next func() (K, V, bool), predicate func(K, V) bool) (K, V, bool) {
	for {
		if k, v, ok := next(); !ok {
			return k, v, false
		} else if ok := predicate(k, v); ok {
			return k, v, true
		}
	}
}

// Firstt returns the first key/value pair that satisfies the condition of the 'predicate' function
func Firstt[K, V any](next func() (K, V, bool), predicate func(K, V) (bool, error)) (K, V, bool, error) {
	for {
		if k, v, ok := next(); !ok {
			return k, v, false, nil
		} else if ok, err := predicate(k, v); err != nil || ok {
			return k, v, ok, err
		}
	}
}

// Convert creates an iterator that applies a transformer to iterable key\values.
func Convert[K, V any, k2, v2 any](next func() (K, V, bool), converter func(K, V) (k2, v2)) ConvertIter[K, V, k2, v2, func(K, V) (k2, v2)] {
	return ConvertIter[K, V, k2, v2, func(K, V) (k2, v2)]{next: next, converter: converter}
}

// Conv creates an iterator that applies a transformer to iterable key\values.
func Conv[K, V any, KOUT, VOUT any](next func() (K, V, bool), converter func(K, V) (KOUT, VOUT, error)) loop.ConvertIter[K, V, KOUT, VOUT] {
	return loop.Conv(loop.From(next), converter)
}

// Filter creates an iterator that checks elements by a filter and returns successful ones
func Filter[K, V any](next func() (K, V, bool), filter func(K, V) bool) FilterIter[K, V] {
	return FilterIter[K, V]{next: next, filter: filter}
}

// Filt creates an iterator that checks elements by a filter and returns successful ones
func Filt[K, V any](next func() (K, V, bool), filter func(K, V) (bool, error)) loop.FiltIter[K, V] {
	return loop.Filt(loop.From(next), filter)
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[K comparable, V, VR any](next func() (K, V, bool), resolver func(bool, K, VR, V) VR) map[K]VR {
	m := map[K]VR{}
	for k, v, ok := next(); ok; k, v, ok = next() {
		exists, ok := m[k]
		m[k] = resolver(ok, k, exists, v)
	}
	return m
}

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[K comparable, V any](next func() (K, V, bool)) map[K]V {
	return ToMapResolv(next, resolv.FirstVal[K, V])
}

// ToSlice collects key\value elements to a slice by iterating over the elements
func ToSlice[K, V, T any](next func() (K, V, bool), converter func(K, V) T) []T {
	s := []T{}
	for key, val, ok := next(); ok; key, val, ok = next() {
		s = append(s, converter(key, val))
	}
	return s
}

func Firs[K, V any](next func() (K, V, bool)) (func() (K, V, bool), K, V, bool) {
	k, v, ok := next()
	return next, k, v, ok
}