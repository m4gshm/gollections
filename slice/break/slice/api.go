// Package slice provides generic functions for slice types
package slice

import (
	"errors"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// ErrBreak is Convert, Filter loops breaker
var ErrBreak = loop.ErrBreak

// ErrIgnore is Convert, Filter element exclude from loop marker
var ErrIgnore = errors.New("Ignore")

// ErrIgnoreAndBreak is loop stopper without the latest element including in the result
var ErrIgnoreAndBreak = errors.New("IgnoreBreak")

// Of is generic slice constructor
func Of[T any](elements ...T) []T { return slice.Of(elements...) }

// OfLoop builds a slice by iterating elements of a source.
// The getNext extracts next element or returns loop break marker or an error.
func OfLoop[S, T any](source S, getNext func(S) (T, error)) ([]T, error) {
	var r []T
	for {
		o, err := getNext(source)
		if err == nil {
			r = append(r, o)
		} else if errors.Is(err, ErrBreak) {
			r = append(r, o)
			return r, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return r, nil
		} else if !errors.Is(err, ErrIgnore) {
			return r, err
		}
	}
}

// Generate builds a slice by an generator function.
// The generator returns an element, or false if the generation is over, or an error.
func Generate[T any](next func() (T, bool, error)) ([]T, error) {
	var r []T
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
		if k, err := keyProducer(e); err == nil {
			initGroup(k, e, groups)
		} else if errors.Is(err, ErrBreak) {
			initGroup(k, e, groups)
			return groups, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return groups, nil
		} else if !errors.Is(err, ErrIgnore) {
			return groups, err
		}
	}
	return groups, nil
}

// GroupInMultiple converts a slice into a map, extracting multiple keys per each element of the slice applying the converter 'keyProducer'
func GroupInMultiple[T any, K comparable, TS ~[]T](elements TS, keysProducer func(T) ([]K, error)) (map[K]TS, error) {
	groups := map[K]TS{}
	for _, e := range elements {
		if keys, err := keysProducer(e); err == nil {
			initGroups(keys, e, groups)
		} else if errors.Is(err, ErrBreak) {
			initGroups(keys, e, groups)
			return groups, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return groups, nil
		} else if !errors.Is(err, ErrIgnore) {
			return groups, err
		}
	}
	return groups, nil
}

func initGroups[T any, K comparable, TS ~[]T](keys []K, e T, groups map[K]TS) {
	if len(keys) == 0 {
		var key K
		initGroup(key, e, groups)
	} else {
		for _, key := range keys {
			initGroup(key, e, groups)
		}
	}
}

func initGroup[T any, K comparable, TS ~[]T](key K, e T, groups map[K]TS) {
	groups[key] = append(groups[key], e)
}

// Convert creates a slice consisting of the transformed elements using the converter 'by'
func Convert[FS ~[]From, From, To any](elements FS, by func(From) (To, error)) ([]To, error) {
	var result []To
	for _, e := range elements {
		if c, err := by(e); err == nil {
			result = append(result, c)
		} else if errors.Is(err, ErrBreak) {
			result = append(result, c)
			return result, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return result, nil
		} else if !errors.Is(err, ErrIgnore) {
			return result, err
		}
	}
	return result, nil
}

// ConvertIndexed creates a slice consisting of the transformed elements using the converter 'by' which additionally applies the index of the element being converted
func ConvertIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, error)) ([]To, error) {
	var result []To
	for i, e := range elements {
		if c, err := by(i, e); err == nil {
			result = append(result, c)
		} else if errors.Is(err, ErrBreak) {
			result = append(result, c)
			return result, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return result, nil
		} else if !errors.Is(err, ErrIgnore) {
			return result, err
		}
	}
	return result, nil
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[FS ~[]From, From, To any](elements FS, by func(From) ([]To, error)) ([]To, error) {
	var result []To
	for _, e := range elements {
		if f, err := by(e); err == nil {
			result = append(result, f...)
		} else if errors.Is(err, ErrBreak) {
			result = append(result, f...)
			return result, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return result, nil
		} else if !errors.Is(err, ErrIgnore) {
			return result, err
		}
	}
	return result, nil
}

// Filter creates a slice containing only the filtered elements
func Filter[TS ~[]T, T any](elements TS, filter func(T) error) ([]T, error) {
	var result []T
	for _, e := range elements {
		if err := filter(e); err == nil {
			result = append(result, e)
		} else if err != nil {
			if errors.Is(err, ErrBreak) {
				result = append(result, e)
				return result, nil
			} else if errors.Is(err, ErrIgnoreAndBreak) {
				return result, nil
			} else if !errors.Is(err, ErrIgnore) {
				return result, err
			}
		}
	}
	return result, nil
}

// Reduce reduces elements to an one
func Reduce[TS ~[]T, T any](elements TS, by func(T, T) (T, error)) (T, error) {
	var result T
	for i, v := range elements {
		if i == 0 {
			result = v
		} else if r, err := by(result, v); err == nil {
			result = r
		} else if errors.Is(err, ErrBreak) {
			result = r
			return result, nil
		} else if errors.Is(err, ErrIgnoreAndBreak) {
			return result, nil
		} else if !errors.Is(err, ErrIgnore) {
			return result, err
		}
	}
	return result, nil
}

// First returns the first element that satisfies requirements of the predicate 'filter'
func First[TS ~[]T, T any](elements TS, by func(T) (bool, error)) (T, bool, error) {
	for _, e := range elements {
		if ok, err := by(e); err != nil {
			var no T
			if errors.Is(err, ErrBreak) || errors.Is(err, ErrIgnoreAndBreak)  {
				return no, false, nil
			} else if !errors.Is(err, ErrIgnore) {
				return no, false, err
			}			
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
			if errors.Is(err, ErrBreak) || errors.Is(err, ErrIgnoreAndBreak)  {
				return no, false, nil
			} else if !errors.Is(err, ErrIgnore) {
				return no, false, err
			}			
		} else if ok {
			return e, true, nil
		}
	}
	var no T
	return no, false, nil
}
