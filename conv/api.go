package conv

import "reflect"

//Converter convert From -> To
type Converter[From, To any] func(From) To

//To helper for Map, Flatt
func To[T any](value T) T { return value }

//AsIs helper for Map, Flatt
func AsIs[T any](value T) T { return value }

//And apply two converters in order
func And[I, O, N any](first Converter[I, O], second Converter[O, N]) Converter[I, N] {
	return func(i I) N { return second(first(i)) }
}

//Or applies first Converter, applies second Converter if the first returns zero
func Or[I, O any](first Converter[I, O], second Converter[I, O]) Converter[I, O] {
	return func(i I) O {
		c := first(i)
		if reflect.ValueOf(c).IsZero() {
			return second(i)
		}
		return c
	}
}

type BinaryOp[T any] func(T, T) T