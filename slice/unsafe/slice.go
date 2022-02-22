package slice

import (
	"reflect"
	"unsafe"
)

//GetTypeSize retrieves size of a type T
func GetTypeSize[T any]() uintptr {
	var t T
	return unsafe.Sizeof(t)
}

//GetArrayPointer retrieves the pointer of a slice underlying array
func GetArrayPointer[T any](elements []T) unsafe.Pointer {
	return unsafe.Pointer(GetSliceHeader(elements).Data)
}

//GetArrayElem returns an element by index from the array referenced by an unsafe pointer
func GetArrayElem[T any](array unsafe.Pointer, index int, elemSize uintptr) T {
	return *(*T)(GetArrayElemRef(array, index, elemSize))
}

func GetArrayElemRef(array unsafe.Pointer, index int, elemSize uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(array) + uintptr(index)*elemSize)
}

//GetLen returns the length of elements
func GetLen[T any](elements []T) int {
	return GetSliceHeader(elements).Len
}

//GetSliceHeader retrieves the SliceHeader of elements
func GetSliceHeader[T any](elements []T) *reflect.SliceHeader {
	return GetSliceHeaderByRef(unsafe.Pointer(&elements))
}

func GetSliceHeaderByRef(elements unsafe.Pointer) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(elements)
}
