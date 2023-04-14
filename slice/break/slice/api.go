// Package slice provides generic functions for slice types
package slice

import (
	"errors"

	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/op"
)

// ErrBreak is Filter, Check breaker
var ErrBreak = it.ErrBreak

// OfLoop builds a slice by iterating elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) ([]T, error) {
	r := []T{}
	for hasNext(source) {
		o, err := getNext(source)
		if err != nil {
			return r, err
		}
		r = append(r, o)
	}
	return r, nil
}

// Generate builds a slice by an generator function.
// The generator returns an element, or false if the generation is over, or an error.
func Generate[T any](next func() (T, bool, error)) ([]T, error) {
	r := []T{}
	for {
		e, ok, err := next()
		if err != nil || !ok {
			return r, err
		}
		r = append(r, e)
	}
}

// DeepClone copies slice elements using a copier function and returns them as a new slice
func DeepClone[TS ~[]T, T any](elements TS, copier func(T) (T, error)) (TS, error) {
	return Convert(elements, copier)
}

// Delete removes an element by index from the slice 'elements'
func Delete[TS ~[]T, T any](index int, elements TS) (TS, error) {
	return append(elements[0:index], elements[index+1:]...), nil
}

// Group converts a slice into a map, extracting a key for each element of the slice applying the converter 'keyProducer'
func Group[T any, K comparable, TS ~[]T](elements TS, keyProducer func(T) (K, error)) (map[K]TS, error) {
	groups := map[K]TS{}
	for _, e := range elements {
		if k, err := keyProducer(e); err != nil {
			return groups, checkBreak(err)
		} else {
			initGroup(k, e, groups)
		}
	}
	return groups, nil
}

// GroupInMultiple converts a slice into a map, extracting multiple keys per each element of the slice applying the converter 'keyProducer'
func GroupInMultiple[T any, K comparable, TS ~[]T](elements TS, keysProducer func(T) ([]K, error)) (map[K]TS, error) {
	groups := map[K]TS{}
	for _, e := range elements {
		if keys, err := keysProducer(e); err != nil {
			return groups, checkBreak(err)
		} else if len(keys) == 0 {
			var key K
			initGroup(key, e, groups)
		} else {
			for _, key := range keys {
				initGroup(key, e, groups)
			}
		}
	}
	return groups, nil
}

func initGroup[T any, K comparable, TS ~[]T](key K, e T, groups map[K]TS) {
	group := groups[key]
	if group == nil {
		group = make([]T, 0)
	}
	groups[key] = append(group, e)
}

// Convert creates a slice consisting of the transformed elements using the converter 'by'
func Convert[FS ~[]From, From, To any](elements FS, by func(From) (To, error)) ([]To, error) {
	result := make([]To, len(elements))
	for i, e := range elements {
		if c, err := by(e); err != nil {
			return result, checkBreak(err)
		} else {
			result[i] = c
		}
	}
	return result, nil
}

// ConvertIndexed creates a slice consisting of the transformed elements using the converter 'by' which additionally applies the index of the element being converted
func ConvertIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, error)) ([]To, error) {
	result := make([]To, len(elements))
	for i, e := range elements {
		if c, err := by(i, e); err != nil {
			return result, checkBreak(err)
		} else {
			result[i] = c
		}
	}
	return result, nil
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[FS ~[]From, From, To any](elements FS, by func(From) ([]To, error)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if f, err := by(e); err != nil {
			return result, checkBreak(err)
		} else {
			result = append(result, f...)
		}
	}
	return result, nil
}

// Filter creates a slice containing only the filtered elements
func Filter[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) ([]T, error) {
	result := make([]T, 0)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return result, checkBreak(err)
		} else if ok {
			result = append(result, e)
		}
	}
	return result, nil
}

// Reduce reduces elements to an one
func Reduce[TS ~[]T, T any](elements TS, by func(T, T) (T, error)) (result T, err error) {
	for i, v := range elements {
		if i == 0 {
			result = v
		} else if result, err = by(result, v); err != nil {
			return result, checkBreak(err)
		}
	}
	return result, nil
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[TS ~[]T, T any](elements TS, by func(T) (bool, error)) (T, bool, error) {
	for _, e := range elements {
		if ok, err := by(e); err != nil {
			var no T
			return no, false, checkBreak(err)
		} else if ok {
			return e, true, nil
		}
	}
	var no T
	return no, false, nil
}

// Last returns the latest element that satisfies requirements of the predicate 'filter'
func Last[TS ~[]T, T any](elements TS, by func(T) (bool, error)) (T, bool, error) {
	for i := len(elements) - 1; i >= 0; i-- {
		e := elements[i]
		if ok, err := by(e); err != nil {
			var no T
			return no, false, checkBreak(err)
		} else if ok {
			return e, true, nil
		}
	}
	var no T
	return no, false, nil
}

func checkBreak(err error) error { return op.IfElse(errors.Is(err, ErrBreak), nil, err) }
