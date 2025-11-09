// Package notsafe provides unsafe helper functions
package notsafe

import (
	"unsafe"
)

// GetTypeSize retrieves size of a type T
func GetTypeSize[T any]() uintptr {
	var t T
	return unsafe.Sizeof(t)
}

// InHeap is generic link to runtime.inheap
func InHeap[T any](t *T) bool {
	return inheap(uintptr(unsafe.Pointer(t)))
}

//go:linkname inheap runtime.inheap
func inheap(p uintptr) bool
