package slice

import (
	"reflect"
)

//Of slice constructor
func Of[T any](values ...T) []T { return values }

//Map Checks items by Predicate filters, applies Converter and accumulate to result slice
func Map[From, To any](items []From, by Converter[From, To], filters ...Predicate[From]) []To {
	result := make([]To, 0)
	for _, v := range items {
		if IsFit[From](v, filters...) {
			result = append(result, by(v))
		}
	}
	return result
}

//Flatt extracts embedded slices of items by Flatter, checks extracted slice values by Predicate filters
//and accumulate to result slice
func Flatt[From, To any](items []From, by Flatter[From, To], filters ...Predicate[To]) []To {
	out := make([]To, 0)
	for _, v := range items {
		flatted := by(v)
		if len(filters) == 0 {
			out = append(out, flatted...)
		} else {
			for _, f := range flatted {
				if IsFit(f, filters...) {
					out = append(out, f)
				}
			}
		}
	}
	return out
}

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

//IsFit apply predicates
func IsFit[T any](v T, predicates ...Predicate[T]) bool {
	fit := true
	for i := 0; fit && i < len(predicates); i++ {
		fit = predicates[i](v)
	}
	return fit
}

//Nil Predicate
func Nil[T any](t T) bool {
	v := reflect.ValueOf(t)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	}
	return false
}

//NotNil Predicate
func NotNil[T any](t T) bool {
	return !Nil(t)
}

//Converter convert From -> To
type Converter[From, To any] func(From) To

//Flatter extracts slice of To
type Flatter[From, To any] func(From) []To

//Predicate tests value (converts to true or false)
type Predicate[T any] Converter[T, bool]
