package check

import (
	"reflect"

	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/typ"
)

//Predicate tests value (converts to true or false)
type Predicate[T any] conv.Converter[T, bool]

//Not invert predicate
func Not[T any](p Predicate[T]) Predicate[T] {
	return func(v T) bool { return !p(v) }
}

//Nil checker.
func Nil[T any](val *T) bool {
	return val == nil
}

//NotNil checker.
func NotNil[T any](val *T) bool {
	return !Nil(val)
}

func Zero[T any](val T) bool {
	return reflect.ValueOf(val).IsZero()
}

func Empty[T typ.Array](val T) bool {
	return len(val) == 0
}

func EmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) == 0
}

func And[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) bool { return p1(v) && p2(v) }
}

//Union reduce predicates to an one
func Union[T any](predicates []Predicate[T]) Predicate[T] {
	l := len(predicates)
	if l == 0 {
		return func(_ T) bool { return false }
	} else if l == 1 {
		return predicates[0]
	} else if l == 2 {
		return And(predicates[0], predicates[1])
	}
	return func(v T) bool {
		for i := 0; i < len(predicates); i++ {
			if !predicates[i](v) {
				return false
			}
		}
		return true
	}
}

func Always[T any](v bool) func(T) bool {
	return func(_ T) bool { return v }
}

func Never[T any](v bool) func(T) bool {
	return func(_ T) bool { return !v }
}
