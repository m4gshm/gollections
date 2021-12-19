package slice

import (
	"reflect"
)

func Of[T any](values ...T) []T { return values }

func AsIs[T any](value T) T { return value }

func And[I, O, N any](f Converter[I, O], s Converter[O, N]) Converter[I, N] {
	return func(i I) N {
		return s(f(i))
	}
}

func Or[I, O any](f Converter[I, O], s Converter[I, O]) Converter[I, O] {
	return func(i I) O {
		c := f(i)
		if reflect.ValueOf(c).IsZero() {
			return s(i)
		}
		return c
	}
}

func Convert[From, To any](values []From, by Converter[From, To]) []To {
	out := make([]To, len(values))
	for i, v := range values {
		o := by(v)
		out[i] = o
	}
	return out
}

func Spread[I, O any](values []I, by Spreader[I, O]) []O {
	result := make([]O, 0)
	for _, c := range values {
		result = append(result, by(c)...)
	}
	return result
}

type Converter[From, To any] func(v From) To
type Spreader[From, To any] func(v From) []To
