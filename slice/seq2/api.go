// Package seq2 provides helpers for rangefunc feature introduced in go 1.22.
package seq2

// Filtered creates a rangefunc that iterates only those elements for which the 'filter' function returns true.
func Filtered[TS ~[]T, T any](elements TS, filter func(T) bool) func(func(int, T) bool) {
	return func(consumer func(int, T) bool) {
		for i, e := range elements {
			if filter(e) {
				if !consumer(i, e) {
					return
				}
			}
		}
	}
}

// Converted creates a rangefunc that applies the 'converter' function to each iterable element.
func Converted[FS ~[]From, From, To any](elements FS, converter func(From) To) func(func(int, To) bool) {
	return func(consumer func(int, To) bool) {
		for i, e := range elements {
			if !consumer(i, converter(e)) {
				return
			}
		}
	}
}
