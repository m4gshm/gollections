package slice

import (
	"reflect"
)

func Of[T any](values ...T) []T { return values }

func AsIs[T any](value T) T { return value }

func And[I, O, N any](first Converter[I, O], second Converter[O, N]) Converter[I, N] {
	return func(i I) N { return second(first(i)) }
}

func Or[I, O any](first Converter[I, O], second Converter[I, O]) Converter[I, O] {
	return func(i I) O {
		c := first(i)
		if reflect.ValueOf(c).IsZero() {
			return second(i)
		}
		return c
	}
}

func Map[From, To any](values []From, by Converter[From, To], filters ...Predicate[From]) []To {
	out := make([]To, 0)
	for _, v := range values {
		if IsFit(v, filters...) {
			c := by(v)
			out = append(out, c)
		}
	}
	return out
}

func Flatt[From, To any](values []From, by Flatter[From, To], filters ...Predicate[To]) []To {
	out := make([]To, 0)
	for _, v := range values {
		flatted := by(v)
		if len(filters) == 0 {
			out = append(out, flatted...)
		} else {
			for _, ss := range flatted {
				if IsFit(ss, filters...) {
					out = append(out, ss)
				}
			}
		}
	}
	return out
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

type Converter[From, To any] func(v From) To
type Flatter[From, To any] func(v From) []To
type Predicate[T any] func(v T) bool
