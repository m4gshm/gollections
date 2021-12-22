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

//IsFit apply predicates
func IsFit[T any](v T, predicates ...Predicate[T]) bool {
	fit := true
	for i := 0; fit && i < len(predicates); i++ {
		fit = predicates[i](v)
	}
	return fit
}
