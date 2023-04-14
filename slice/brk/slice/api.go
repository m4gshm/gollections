// Package slice provides generic functions for slice types
package slice

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
func Group[T any, K comparable, TS ~[]T](elements TS, keyProducer func(T) K) map[K]TS {
	groups := map[K]TS{}
	for _, e := range elements {
		initGroup(keyProducer(e), e, groups)
	}
	return groups
}

// GroupInMultiple converts a slice into a map, extracting multiple keys per each element of the slice applying the converter 'keyProducer'
func GroupInMultiple[T any, K comparable, TS ~[]T](elements TS, keysProducer func(T) []K) map[K]TS {
	groups := map[K]TS{}
	for _, e := range elements {
		if keys := keysProducer(e); len(keys) == 0 {
			var key K
			initGroup(key, e, groups)
		} else {
			for _, key := range keys {
				initGroup(key, e, groups)
			}
		}
	}
	return groups
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
			return result, err
		} else {
			result[i] = c
		}
	}
	return result, nil
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), by func(From) (To, error)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return result, err
		} else if ok {
			if c, err := by(e); err != nil {
				return result, err
			} else {
				result = append(result, c)
			}
		}
	}
	return result, nil
}

// FilterAndConvert additionally filters 'To' elements
func ConvertAndFilter[FS ~[]From, From, To any](elements FS, by func(From) (To, error), filter func(To) (bool, error)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if r, err := by(e); err != nil {
			return result, err
		} else if ok, err := filter(r); err != nil {
			return result, err
		} else if ok {
			result = append(result, r)
		}
	}
	return result, nil
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), by func(From) (To, error), filterConverted func(To) (bool, error)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return result, err
		} else if ok {
			if r, err := by(e); err != nil {
				return result, err
			} else if ok, err := filterConverted(r); err != nil {
				return result, err
			} else if ok {
				result = append(result, r)
			}
		}
	}
	return result, nil
}

// ConvertIndexed creates a slice consisting of the transformed elements using the converter 'by' which additionally applies the index of the element being converted
func ConvertIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) To) ([]To, error) {
	result := make([]To, len(elements))
	for i, e := range elements {
		result[i] = by(i, e)
	}
	return result, nil
}

// FilterAndConvertIndexed additionally filters 'From' elements
func FilterAndConvertIndexed[FS ~[]From, From, To any](elements FS, filter func(index int, from From) bool, converter func(index int, from From) To) ([]To, error) {
	result := make([]To, 0)
	for i, e := range elements {
		if filter(i, e) {
			result = append(result, converter(i, e))
		}
	}
	return result, nil
}

// ConvertCheck is similar to ConvertFit, but it checks and transforms elements together
func ConvertCheck[FS ~[]From, From, To any](elements FS, by func(from From) (To, bool)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if to, ok := by(e); ok {
			result = append(result, to)
		}
	}
	return result, nil
}

// ConvertCheckIndexed additionally filters 'From' elements
func ConvertCheckIndexed[FS ~[]From, From, To any](elements FS, by func(index int, from From) (To, bool)) ([]To, error) {
	result := make([]To, 0)
	for i, e := range elements {
		if to, ok := by(i, e); ok {
			result = append(result, to)
		}
	}
	return result, nil
}

// Flatt unfolds the n-dimensional slice into a n-1 dimensional slice
func Flatt[FS ~[]From, From, To any](elements FS, by func(From) []To) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		result = append(result, by(e)...)

	}
	return result, nil
}

// FilerAndFlatt additionally filters 'From' elements.
func FilerAndFlatt[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), by func(From) []To) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return result, err
		} else if ok {
			result = append(result, by(e)...)
		}
	}
	return result, nil
}

// FlattAndFiler unfolds the n-dimensional slice into a n-1 dimensional slice with additinal filtering of 'To' elements.
func FlattAndFiler[FS ~[]From, From, To any](elements FS, by func(From) []To, filter func(To) (bool, error)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		for _, to := range by(e) {
			if ok, err := filter(to); err != nil {
				return result, err
			} else if ok {
				result = append(result, to)
			}
		}
	}
	return result, nil
}

// FilterFlattFilter unfolds the n-dimensional slice 'elements' into a n-1 dimensional slice with additinal filtering of 'From' and 'To' elements.
func FilterFlattFilter[FS ~[]From, From, To any](elements FS, filterFrom func(From) (bool, error), by func(From) []To, filterTo func(To) (bool, error)) ([]To, error) {
	result := make([]To, 0)
	for _, e := range elements {
		if ok, err := filterFrom(e); err != nil {
			return result, err
		} else if ok {
			for _, to := range by(e) {
				if ok, err := filterTo(to); err != nil {
					return result, err
				} else if ok {
					result = append(result, to)
				}
			}
		}
	}
	return result, nil
}

// Filter creates a slice containing only the filtered elements
func Filter[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) ([]T, error) {
	result := make([]T, 0)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return result, err
		} else if ok {
			result = append(result, e)
		}
	}
	return result, nil
}

// FilterBrk creates a slice containing only the filtered elements or aborts filtering by an error
func FilterBrk[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) ([]T, error) {
	result := make([]T, 0)
	for _, e := range elements {
		if ok, err := filter(e); err != nil {
			return nil, err
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
			return no, false, err
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
			return no, false, err
		} else if ok {
			return e, true, nil
		}
	}
	var no T
	return no, false, nil
}
