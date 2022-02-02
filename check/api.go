package check

import (
	"reflect"

	"github.com/m4gshm/gollections/c"
)

//Not inverts a predicate.
func Not[T any](p c.Predicate[T]) c.Predicate[T] {
	return func(v T) bool { return !p(v) }
}

//Nil checks a reference for nil value.
func Nil[T any](val *T) bool {
	return val == nil
}

//NotNil checks a reference for no nil value.
func NotNil[T any](val *T) bool {
	return !Nil(val)
}

//Zero checks that a value is zero.
func Zero[T any](val T) bool {
	return reflect.ValueOf(val).IsZero()
}

//Empty checks that a slice is empty.
func Empty[T Slice](val T) bool {
	return len(val) == 0
}

//EmptyMap checks that a slice is ampty.
func EmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) == 0
}

//And makes a conjunction of two predicates
func And[T any](p1, p2 c.Predicate[T]) c.Predicate[T] {
	return func(v T) bool { return p1(v) && p2(v) }
}

//Union applies And to predicates
func Union[T any](predicates []c.Predicate[T]) c.Predicate[T] {
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

//Always returns v every time
func Always[T any](v bool) c.Predicate[T] {
	return func(_ T) bool { return v }
}

//Always returns the negative of v every time
func Never[T any](v bool) c.Predicate[T] {
	return func(_ T) bool { return !v }
}

type Slice interface {
	~[]any | ~[]uintptr |
		~[]int | ~[]int8 | []int16 | []int32 | []int64 |
		~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 |
		~[]float32 | ~[]float64 |
		~[]complex64 | ~[]complex128 |
		~[]string | ~string
}
