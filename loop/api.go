package loop

import "errors"

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = errors.New("Break")

// For applies the 'walker' function for the elements retrieved by the 'next' function. Return the c.ErrBreak to stop
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

// ForEach applies the 'walker' function to the elements retrieved by the 'next' function
func ForEach[T any](next func() (T, bool), walker func(T)) {
	for v, ok := next(); ok; v, ok = next() {
		walker(v)
	}
}

// ForEachFiltered applies the 'walker' function to the elements retrieved by the 'next' function that satisfy the condition of the 'predicate' function
func ForEachFiltered[T any](next func() (T, bool), walker func(T), predicate func(T) bool) {
	for v, ok := next(); ok && predicate(v); v, ok = next() {
		walker(v)
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool), predicate func(T) bool) (v T, ok bool) {
	for one, ok := next(); ok; one, ok = next() {
		if predicate(one) {
			return one, true
		}
	}
	return v, ok
}

// Track applies the 'tracker' function to position/element pairs retrieved by the 'next' function. Return the c.ErrBreak to stop tracking..
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

// TrackEach applies the 'tracker' function to position/element pairs retrieved by the 'next' function
func TrackEach[I, T any](next func() (I, T, bool), tracker func(I, T)) {
	for p, v, ok := next(); ok; p, v, ok = next() {
		tracker(p, v)
	}
}

// ToSlice collects the elements retrieved by the 'next' function into a slice
func ToSlice[T any](next func() (T, bool)) []T {
	var s []T
	for v, ok := next(); ok; v, ok = next() {
		s = append(s, v)
	}
	return s
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function
func Reduce[T any](next func() (T, bool), merger func(T, T) T) (result T) {
	if v, ok := next(); ok {
		result = v
	} else {
		return result
	}
	for v, ok := next(); ok; v, ok = next() {
		result = merger(result, v)
	}
	return result
}

// ReduceKV reduces the key/value pairs retrieved by the 'next' function into an one pair using the 'merge' function
func ReduceKV[K, V any](next func() (K, V, bool), merge func(K, V, K, V) (K, V)) (rk K, rv V) {
	if k, v, ok := next(); ok {
		rk, rv = k, v
	} else {
		return rk, rv
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		rk, rv = merge(rk, rv, k, v)
	}
	return rk, rv
}
