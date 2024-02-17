package over

import "github.com/m4gshm/gollections/c"

func Filtered[T any](all c.RangeFunc[T], filter func(T) bool) c.RangeFunc[T] {
	return func(consumerFiltered func(T) bool) {
		all(func(e T) bool {
			if filter(e) {
				return consumerFiltered(e)
			}
			return true
		})
	}
}

func Converted[From, To any](all c.RangeFunc[From], converter func(From) To) c.RangeFunc[To] {
	return func(consumerTo func(To) bool) {
		all(func(from From) bool {
			return consumerTo(converter(from))
		})
	}
}
