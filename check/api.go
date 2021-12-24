package check

import (
	"reflect"

	"github.com/m4gshm/container/conv"
)

//Predicate tests value (converts to true or false)
type Predicate[T any] conv.Converter[T, bool]

//Nil checker. Safe for non-nullable types
func Nil[T any](val T) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.Invalid:
		return reflect.TypeOf(val) == nil
	}
	return false
}

//NotNil checker. Safe for non-nullable types
func NotNil[T any](val T) bool {
	return !Nil(val)
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
