package over

import "github.com/m4gshm/gollections/c"

func Filtered[T any](all c.RangeFunc[T], filter func(T) bool) c.RangeFunc[T] {
	return func(yieldFiltered func(T) bool) {
		all(func(e T) bool {
			if filter(e) {
				return yieldFiltered(e)
			}
			return true
		})
	}
}

func Converted[From, To any](all c.RangeFunc[From], converter func(From) To) c.RangeFunc[To] {
	return func(yieldTo func(To) bool) {
		all(func(from From) bool {
			return yieldTo(converter(from))
		})
	}
}
