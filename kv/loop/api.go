// Package loop provides helpers for loop operation over key/value pairs and iterator implementations
//
// Deprecated: use the [github.com/m4gshm/gollections/seq], [github.com/m4gshm/gollections/seqe], [github.com/m4gshm/gollections/seq2] packages API instead.
package loop

import (
	"github.com/m4gshm/gollections/break/kv/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/map_/resolv"
)

// New makes a seq from an abstract source
func New[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V)) Loop[K, V] {
	return func() (k K, v V, ok bool) {
		if hasNext(source) {
			k, v = getNext(source)
			return k, v, true
		}
		return k, v, false
	}
}

// All is an adapter for the next function for iterating by `for ... range`.
func All[K, V any](next func() (K, V, bool), consumer func(K, V) bool) {
	for k, v, ok := next(); ok && consumer(k, v); k, v, ok = next() {
	}
}

// Track applies the 'consumer' function to position/element pairs retrieved by the 'next' function until the consumer returns the c.Break to stop.
func Track[I, T any](next func() (I, T, bool), consumer func(I, T) error) error {
	if next == nil {
		return nil
	}
	for p, v, ok := next(); ok; p, v, ok = next() {
		if err := consumer(p, v); err == c.Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEach applies the 'consumer' function to position/element pairs retrieved by the 'next' function
func TrackEach[I, T any](next func() (I, T, bool), consumer func(I, T)) {
	if next == nil {
		return
	}
	for p, v, ok := next(); ok; p, v, ok = next() {
		consumer(p, v)
	}
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](next func() (K, V, bool)) map[K][]V {
	return MapResolv(next, resolv.Slice[K, V])
}

// Reduce reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero values of 'K', 'V' types are returned.
func Reduce[K, V any](next func() (K, V, bool), merge func(K, K, V, V) (K, V)) (rk K, rv V) {
	rk, rv, _ = ReduceOK(next, merge)
	return rk, rv
}

// ReduceOK reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceOK[K, V any](next func() (K, V, bool), merge func(K, K, V, V) (K, V)) (rk K, rv V, ok bool) {
	if next == nil {
		return rk, rv, false
	}
	k, v, ok := next()
	if !ok {
		return k, v, false
	}
	rk, rv = k, v
	for k, v, ok := next(); ok; k, v, ok = next() {
		rk, rv = merge(rk, k, rv, v)
	}
	return rk, rv, true
}

// Reducee reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero values of 'K', 'V' types are returned.
func Reducee[K, V any](next func() (K, V, bool), merge func(K, K, V, V) (K, V, error)) (rk K, rv V, err error) {
	rk, rv, _, err = ReduceeOK(next, merge)
	return rk, rv, err
}

// ReduceeOK reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceeOK[K, V any](next func() (K, V, bool), merge func(K, K, V, V) (K, V, error)) (rk K, rv V, ok bool, err error) {
	if next == nil {
		return rk, rv, false, nil
	}
	k, v, ok := next()
	if !ok {
		return rk, rv, false, nil
	}
	rk, rv = k, v
	for {
		if k, v, ok := next(); !ok {
			return rk, rv, true, nil
		} else if rk, rv, err = merge(rk, k, rv, v); err != nil {
			return rk, rv, true, err
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

// Convert creates a seq that applies the 'converter' function to iterable key\values.
func Convert[K, V any, KOUT, VOUT any](next func() (K, V, bool), converter func(K, V) (KOUT, VOUT)) Loop[KOUT, VOUT] {
	if next == nil {
		return nil
	}
	return func() (k2 KOUT, v2 VOUT, ok bool) {
		if k, v, ok := next(); ok {
			k2, v2 = converter(k, v)
			return k2, v2, true
		}
		return k2, v2, false
	}
}

// Conv creates a seq that applies the 'converter' function to iterable key\values.
func Conv[K, V any, KOUT, VOUT any](next func() (K, V, bool), converter func(K, V) (KOUT, VOUT, error)) loop.Loop[KOUT, VOUT] {
	return loop.Conv(loop.From(next), converter)
}

// Filter creates a seq that checks elements by the 'filter' function and returns successful ones.
func Filter[K, V any](next func() (K, V, bool), filter func(K, V) bool) Loop[K, V] {
	if next == nil {
		return nil
	}
	return func() (K, V, bool) {
		return First(next, filter)
	}
}

// Filt creates a seq that checks elements by the 'filter' function and returns successful ones.
func Filt[K, V any](next func() (K, V, bool), filter func(K, V) (bool, error)) loop.Loop[K, V] {
	return loop.Filt(loop.From(next), filter)
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values
func MapResolv[K comparable, V, VR any](next func() (K, V, bool), resolver func(bool, K, VR, V) VR) map[K]VR {
	if next == nil {
		return nil
	}
	m := map[K]VR{}
	for k, v, ok := next(); ok; k, v, ok = next() {
		exists, ok := m[k]
		m[k] = resolver(ok, k, exists, v)
	}
	return m
}

// Map collects key\value elements into a new map by iterating over the elements
func Map[K comparable, V any](next func() (K, V, bool)) map[K]V {
	return MapResolv(next, resolv.First[K, V])
}

// Slice collects key\value elements to a slice by iterating over the elements
func Slice[K, V, T any](next func() (K, V, bool), converter func(K, V) T) []T {
	if next == nil {
		return nil
	}
	s := []T{}
	for key, val, ok := next(); ok; key, val, ok = next() {
		s = append(s, converter(key, val))
	}
	return s
}

// Crank rertieves next key\value from the 'next' function, returns the function, element, successfully flag.
func Crank[K, V any](next func() (K, V, bool)) (n Loop[K, V], k K, v V, ok bool) {
	if next != nil {
		k, v, ok = next()
	}
	return next, k, v, ok
}
