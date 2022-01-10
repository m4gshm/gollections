package check

import (
	"reflect"

	"github.com/m4gshm/gollections/typ"
)

//Not invert predicate
func Not[T any](p typ.Predicate[T]) typ.Predicate[T] {
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

func Empty[T Array](val T) bool {
	return len(val) == 0
}

func EmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) == 0
}

func And[T any](p1, p2 typ.Predicate[T]) typ.Predicate[T] {
	return func(v T) bool { return p1(v) && p2(v) }
}

//Union reduce predicates to an one
func Union[T any](predicates []typ.Predicate[T]) typ.Predicate[T] {
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

type Array interface {
	~[]any | ~[]uintptr |
		~[]int | ~[]int8 | []int16 | []int32 | []int64 |
		~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 |
		~[]float32 | ~[]float64 |
		~[]complex64 | ~[]complex128 |
		~[]string | ~string
}
