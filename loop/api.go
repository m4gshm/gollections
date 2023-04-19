package loop

import "errors"

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = errors.New("Break")

// For applies a walker to elements retrieved by the 'next' function. To stop walking just return the ErrBreak
func For[T any](next func() (T, bool), walker func(T) error) error {
	for v, ok := next(); ok; v, ok = next() {
		if err := walker(v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies a walker to elements of an Iterator
func ForEach[T any](next func() (T, bool), walker func(T)) {
	for v, ok := next(); ok; v, ok = next() {
		walker(v)
	}
}

// ForEachFiltered applies a walker to elements that satisfy a predicate condition
func ForEachFiltered[T any](next func() (T, bool), walker func(T), filter func(T) bool) {
	for v, ok := next(); ok && filter(v); v, ok = next() {
		walker(v)
	}
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[T any](next func() (T, bool), filter func(T) bool) (v T, ok bool) {
	for one, ok := next(); ok; one, ok = next() {
		if filter(one) {
			return one, true
		}
	}
	return v, ok
}

// Track apply a tracker to position/element pairs retrieved by the 'next' function. To stop traking just return the ErrBreak.
func Track[I, T any](next func() (I, T, bool), tracker func(I, T) error) error {
	for p, v, ok := next(); ok; p, v, ok = next() {
		if err := tracker(p, v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEach applies a tracker for every position/value pairs from a indexed iterator (array, map)
func TrackEach[I, T any](next func() (I, T, bool), tracker func(I, T)) {
	for p, v, ok := next(); ok; p, v, ok = next() {
		tracker(p, v)
	}
}

// ToSlice collects elements retrieved by the 'next' function into a slice
func ToSlice[T any](next func() (T, bool)) []T {
	var s []T
	for v, ok := next(); ok; v, ok = next() {
		s = append(s, v)
	}
	return s
}

// Reduce reduces elements to an one
func Reduce[T any](next func() (T, bool), by func(T, T) T) (result T) {
	if v, ok := next(); ok {
		result = v
	} else {
		return result
	}
	for v, ok := next(); ok; v, ok = next() {
		result = by(result, v)
	}
	return result
}

// ReduceKV reduces key/values elements to an one
func ReduceKV[K, V any](next func() (K, V, bool), by func(K, V, K, V) (K, V)) (rk K, rv V) {
	if k, v, ok := next(); ok {
		rk, rv = k, v
	} else {
		return rk, rv
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		rk, rv = by(rk, rv, k, v)
	}
	return rk, rv
}
