package check

import (
	"reflect"
)

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
