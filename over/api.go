// Package over provides helpers for rangefunc feature introduced in go 1.22.
package over

// Filtered creates a rangefunc that iterates only those elements for which the 'filter' function returns true.
func Filtered[T any](all func(func(T) bool), filter func(T) bool) func(func(T) bool) {
	return func(consumer func(T) bool) {
		all(func(e T) bool {
			if filter(e) {
				return consumer(e)
			}
			return true
		})
	}
}

// Converted creates a rangefunc that applies the 'converter' function to each iterable element.
func Converted[From, To any](all func(func(From) bool), converter func(From) To) func(func(To) bool) {
	return func(consumer func(To) bool) {
		all(func(from From) bool {
			return consumer(converter(from))
		})
	}
}
