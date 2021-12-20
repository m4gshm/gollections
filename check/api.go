package check

import "reflect"

//Nil checker. Safe for non-nullable types
func Nil[T any](t T) bool {
	v := reflect.ValueOf(t)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	}
	return false
}

//NotNil checker. Safe for non-nullable types
func NotNil[T any](t T) bool {
	return !Nil(t)
}
