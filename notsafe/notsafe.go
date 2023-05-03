// Package notsafe provides unsafe helper functions
package notsafe

import (
	"reflect"
	"unsafe"
)

// GetTypeSize retrieves size of a type T
func GetTypeSize[T any]() uintptr {
	var t T
	return unsafe.Sizeof(t)
}

// GetArrayPointer retrieves the pointer of a slice underlying array
func GetArrayPointer[T any](elements []T) unsafe.Pointer {
	return unsafe.Pointer(GetSliceHeader(elements).Data)
}

// GetArrayElem returns an element by index from the array referenced by an unsafe pointer
func GetArrayElem[T any](array unsafe.Pointer, index int, elemSize uintptr) T {
	return *(*T)(GetArrayElemRef(array, index, elemSize))
}

// GetArrayElemRef returns an element's pointer by index from the array referenced by an unsafe pointer
func GetArrayElemRef(array unsafe.Pointer, index int, elemSize uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(array) + uintptr(index)*elemSize)
}

// GetLen returns the length of elements
func GetLen[T any](elements []T) int {
	return GetSliceHeader(elements).Len
}

// GetSliceHeader retrieves the SliceHeader of elements
func GetSliceHeader[T any](elements []T) *reflect.SliceHeader {
	return GetSliceHeaderByRef(unsafe.Pointer(&elements))
}

// GetSliceHeaderByRef retrieves the SliceHeader of elements by an unsafe pointer
func GetSliceHeaderByRef(elements unsafe.Pointer) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(elements)
}

// Noescape prevent escaping of t
//
//go:nosplit
//go:nocheckptr
func Noescape[T any](t *T) *T {
	x := uintptr(unsafe.Pointer(t))
	return (*T)(unsafe.Pointer(x ^ 0)) //nolint
}

// InHeap is generic link to runtime.inheap
func InHeap[T any](t *T) bool {
	return inheap(uintptr(unsafe.Pointer(t)))
}

//go:linkname inheap runtime.inheap
func inheap(p uintptr) bool
