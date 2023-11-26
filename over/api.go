package over

import "github.com/m4gshm/gollections/c"

func Filtered[T any](all c.RangeFunc[T], filter func(T) bool) c.RangeFunc[T] {
	return func(yieldFiltered func(T) bool) {
		yield := func(e T) bool {
			if filter(e) {
				return yieldFiltered(e)
			}
			return true
		}
		all(yield)
	}
}

func Converted[From, To any](all c.RangeFunc[From], converter func(From) To) c.RangeFunc[To] {
	return func(yieldTo func(To) bool) {
		yieldFrom := func(f From) bool {
			to := converter(f)
			return yieldTo(to)
		}
		all(yieldFrom)
	}
}
