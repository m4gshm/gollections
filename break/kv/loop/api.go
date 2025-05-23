// Package loop provides helpers for loop operation over key/value pairs.
//
// Deprecated: use the [github.com/m4gshm/gollections/seq], [github.com/m4gshm/gollections/seqe], [github.com/m4gshm/gollections/seq2] packages API instead.
package loop

import (
	"errors"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/map_/resolv"
)

// New is the mai breakable key/value loop constructor
func New[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) Loop[K, V] {
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

// From wrap the next loop to a breakable loop
func From[K, V any](next func() (K, V, bool)) func() (K, V, bool, error) {
	if next == nil {
		return nil
	}
	return func() (K, V, bool, error) {
		k, v, ok := next()
		return k, v, ok, nil
	}
}

// To transforms a breakable loop to a simple loop.
// The errConsumer is a function that is called when an error occurs.
func To[K, V any](next func() (K, V, bool, error), errConsumer func(error)) func() (K, V, bool) {
	if next == nil {
		return nil
	}
	return func() (K, V, bool) {
		k, v, ok, err := next()
		if err != nil {
			errConsumer(err)
			return k, v, false
		}
		return k, v, ok
	}
}

// Group collects sets of values grouped by keys obtained by passing a key/value loop.
func Group[K comparable, V any](next func() (K, V, bool, error)) (map[K][]V, error) {
	return MapResolv(next, resolv.Slice[K, V])
}

// Reduce reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero values of 'K', 'V' types are returned.
func Reduce[K, V any](next func() (K, V, bool, error), merge func(K, K, V, V) (K, V)) (K, V, error) {
	rk, rv, _, err := ReduceOK(next, merge)
	return rk, rv, err
}

// ReduceOK reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceOK[K, V any](next func() (K, V, bool, error), merge func(K, K, V, V) (K, V)) (rk K, rv V, ok bool, err error) {
	if next == nil {
		return rk, rv, false, nil
	}
	k, v, ok, err := next()
	if err != nil || !ok {
		return k, v, ok, err
	}
	rk, rv = k, v
	for {
		k, v, ok, err := next()
		if err != nil || !ok {
			return rk, rv, true, err
		}
		rk, rv = merge(rk, k, rv, v)
	}
}

// Reducee reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero values of 'K', 'V' types are returned.
func Reducee[K, V any](next func() (K, V, bool, error), merge func(K, K, V, V) (K, V, error)) (K, V, error) {
	rk, rv, _, err := ReduceeOK(next, merge)
	return rk, rv, err
}

// ReduceeOK reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceeOK[K, V any](next func() (K, V, bool, error), merge func(K, K, V, V) (K, V, error)) (rk K, rv V, ok bool, err error) {
	if next == nil {
		return rk, rv, false, nil
	}
	k, v, ok, err := next()
	if err != nil || !ok {
		return rk, rv, ok, err
	}
	rk, rv = k, v
	for {
		if k, v, ok, err := next(); err != nil || !ok {
			return rk, rv, true, err
		} else if rk, rv, err = merge(rk, k, rv, v); err != nil {
			return rk, rv, true, err
		}
	}
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAny[K, V any](next func() (K, V, bool, error), predicate func(K, V) bool) (bool, error) {
	_, _, ok, err := First(next, predicate)
	return ok, err
}

// HasAnyy finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func HasAnyy[K, V any](next func() (K, V, bool, error), predicate func(K, V) (bool, error)) (bool, error) {
	_, _, ok, err := Firstt(next, predicate)
	return ok, err
}

// First returns the first key/value pair that satisfies the condition of the 'predicate' function
func First[K, V any](next func() (K, V, bool, error), predicate func(K, V) bool) (K, V, bool, error) {
	for {
		if k, v, ok, err := next(); err != nil || !ok {
			return k, v, false, err
		} else if ok := predicate(k, v); ok {
			return k, v, true, nil
		}
	}
}

// Firstt returns the first key/value pair that satisfies the condition of the 'predicate' function
func Firstt[K, V any](next func() (K, V, bool, error), predicate func(K, V) (bool, error)) (K, V, bool, error) {
	for {
		if k, v, ok, err := next(); err != nil || !ok {
			return k, v, false, err
		} else if ok, err := predicate(k, v); err != nil || ok {
			return k, v, ok, err
		}
	}
}

// Convert creates a loop that applies the 'converter' function to iterable key\values.
func Convert[K, V any, KOUT, VOUT any](next func() (K, V, bool, error), converter func(K, V) (KOUT, VOUT)) Loop[KOUT, VOUT] {
	if next == nil {
		return nil
	}
	return func() (k2 KOUT, v2 VOUT, ok bool, err error) {
		k, v, ok, err := next()
		if err != nil || !ok {
			return k2, v2, false, err
		}
		k2, v2 = converter(k, v)
		return k2, v2, true, nil
	}
}

// Conv creates a loop that applies the 'converter' function to iterable key\values.
func Conv[K, V any, KOUT, VOUT any](next func() (K, V, bool, error), converter func(K, V) (KOUT, VOUT, error)) Loop[KOUT, VOUT] {
	if next == nil {
		return nil
	}
	return func() (k2 KOUT, v2 VOUT, ok bool, err error) {
		k, v, ok, err := next()
		if err != nil || !ok {
			return k2, v2, false, err
		}
		k2, v2, err = converter(k, v)
		return k2, v2, true, err
	}
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func Filter[K, V any](next func() (K, V, bool, error), filter func(K, V) bool) Loop[K, V] {
	if next == nil {
		return nil
	}
	return func() (K, V, bool, error) {
		return First(next, filter)
	}
}

// Filt creates a loop that checks elements by the 'filter' function and returns successful ones.
func Filt[K, V any](next func() (K, V, bool, error), filter func(K, V) (bool, error)) Loop[K, V] {
	if next == nil {
		return nil
	}
	return func() (K, V, bool, error) {
		k, v, ok, err := Firstt(next, filter)
		return k, v, ok && err == nil, err
	}
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values
func MapResolv[K comparable, V, VR any](next func() (K, V, bool, error), resolver func(bool, K, VR, V) VR) (map[K]VR, error) {
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

// Map collects key\value elements into a new map by iterating over the elements
func Map[K comparable, V any](next func() (K, V, bool, error)) (map[K]V, error) {
	return MapResolv(next, resolv.First[K, V])
}

// Slice collects key\value elements to a slice by iterating over the elements
func Slice[K, V, T any](next func() (K, V, bool, error), converter func(K, V) T) ([]T, error) {
	if next == nil {
		return nil, nil
	}
	s := []T{}
	for {
		key, val, ok, err := next()
		if ok {
			s = append(s, converter(key, val))
		}
		if !ok || err != nil {
			return s, err
		}
	}
}

// Track applies the 'consumer' function to position/element pairs retrieved by the 'next' function until the consumer returns the c.Break to stop.
func Track[K, V any](next func() (K, V, bool, error), consumer func(K, V) error) error {
	if next == nil {
		return nil
	}
	for {
		if p, v, ok, err := next(); err != nil || !ok {
			return err
		} else if err := consumer(p, v); err != nil {
			return brk(err)
		}
	}
}

// Crank rertieves a next element from the 'next' function, returns the function, element, successfully flag.
func Crank[K, V any](next func() (K, V, bool, error)) (n Loop[K, V], k K, v V, ok bool, err error) {
	if next != nil {
		k, v, ok, err = next()
	}
	return next, k, v, ok, err
}

func brk(err error) error {
	if errors.Is(err, c.Break) {
		return nil
	}
	return err
}
