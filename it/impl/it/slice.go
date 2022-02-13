package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

const NoStarted = -1

func NewHead[T any](elements []T) *Iter[T] {
	return NewHeadS(elements, GetTypeSize[T]())
}

func NewHeadS[T any](elements []T, elementSize uintptr) *Iter[T] {
	header := GetSliceHeader(elements)
	return &Iter[T]{
		array:       unsafe.Pointer(header.Data),
		elementSize: elementSize,
		size:        header.Len,
		current:     NoStarted,
	}
}

func NewTail[T any](elements []T) *Iter[T] {
	return NewTailS(elements, GetTypeSize[T]())
}

func NewTailS[T any](elements []T, elementSize uintptr) *Iter[T] {
	var (
		header = GetSliceHeader(elements)
		size   = header.Len
	)
	return &Iter[T]{
		array:       unsafe.Pointer(header.Data),
		elementSize: elementSize,
		size:        size,
		current:     size,
	}
}

type Iter[T any] struct {
	array         unsafe.Pointer
	elementSize   uintptr
	size, current int
}

var _ c.Iterator[any] = (*Iter[any])(nil)

func (i *Iter[T]) HasNext() bool {
	if HasNextBySize(i.size, i.current) {
		i.current++
		return true
	}
	return false
}

func (i *Iter[T]) HasPrev() bool {
	if HasPrevBySize(i.size, i.current) {
		i.current--
		return true
	}
	return false
}

func (i *Iter[T]) Get() T {
	return GetArrayElem[T](i.array, i.current, i.elementSize)
}

func HasNext[T any](elements []T, current int) bool {
	return HasNextBySize(GetLen(elements), current)
}

func HasPrev[T any](elements []T, current int) bool {
	return HasPrevBySize(GetLen(elements), current)
}

func HasNextBySize(size int, current int) bool {
	return HasNextByRange(NoStarted, size-2, current)
}

func HasPrevBySize(size, current int) bool {
	return HasNextByRange(1, size, current)
}

func HasNextByRange(first, last, current int) bool {
	return current >= first && current <= last
}

func Get[T any](elements []T, current int) T {
	if current >= len(elements) {
		var no T
		return no
	} else if current == NoStarted {
		var no T
		return no
	}
	return (elements)[current]
}

func GetArrayPointer[T any](elements []T) unsafe.Pointer {
	return unsafe.Pointer(GetSliceHeader(elements).Data)
}

func GetArrayPointer2[T any](elements []T) uintptr {
	return GetSliceHeader(elements).Data
}

func GetTypeSize[T any]() uintptr {
	var t T
	return unsafe.Sizeof(t)
}

func GetArrayElem[T any](array unsafe.Pointer, index int, elemSize uintptr) T {
	return *(*T)(unsafe.Pointer(uintptr(array) + uintptr(index)*elemSize))
}

func GetLen[T any](elements []T) int {
	return GetSliceHeader(elements).Len
}

func GetSliceHeader[T any](elements []T) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&elements))
}
