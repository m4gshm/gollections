package check

import (
	"github.com/m4gshm/container/conv"
)

//Predicate tests value (converts to true or false)
type Predicate[T any] conv.Converter[T, bool]

//Nil checker.
func Nil[T any](val *T) bool {
	return val == nil
}

//NotNil checker.
func NotNil[T any](val *T) bool {
	return !Nil(val)
}

type array interface {
	[]any | []int | []int8 | []int16 | []int32 | []int64 | []uint | []uint8 | []uint16 | []uint32 | []uint64 | []float32 | []float64 | []uintptr | []complex64 | []complex128 | []string | string
}

func Empty[T array](val T) bool {
	return len(val) == 0
}

func EmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) == 0
}

func Not[T any](p Predicate[T]) Predicate[T] {
	return func(v T) bool { return !p(v) }
}

//Union reduce predicates to an one
func Union[T any](predicates []Predicate[T]) Predicate[T] {
	l := len(predicates)
	if l == 0 {
		return func(_ T) bool { return false }
	} else if l == 1 {
		return predicates[0]
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
	return func(_ T) bool { return true }
}

func Never[T any](v bool) func(T) bool {
	return func(_ T) bool { return true }
}
