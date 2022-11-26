package conv

import (
	"reflect"

	"github.com/m4gshm/gollections/c"
)

// To helper for Map, Flatt.
func To[T any](value T) T { return value }

// AsIs helper for Map, Flatt.
func AsIs[T any](value T) T { return value }

// And apply two converters in order.
func And[I, O, N any](first c.Converter[I, O], second c.Converter[O, N]) c.Converter[I, N] {
	return func(i I) N { return second(first(i)) }
}

// Or applies first Converter, applies second Converter if the first returns zero.
func Or[I, O any](first c.Converter[I, O], second c.Converter[I, O]) c.Converter[I, O] {
	return func(i I) O {
		c := first(i)
		if reflect.ValueOf(c).IsZero() {
			return second(i)
		}
		return c
	}
}
