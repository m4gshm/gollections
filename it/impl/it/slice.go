package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

const NoStarted = -1

func NewHead[T any](elements []T) *Iter[T] {
	return &Iter[T]{
		array:            GetArrayPointer(elements),
		arraySize:        GetLen(elements),
		arrayElementSize: GetArrayElementSize(elements),
		current:          NoStarted,
	}
}

func NewTail[T any](elements []T) *Iter[T] {
	l := GetLen(elements)
	return &Iter[T]{
		array:            GetArrayPointer(elements),
		arraySize:        l,
		arrayElementSize: GetArrayElementSize(elements),
		current:          l,
	}
}

type Iter[T any] struct {
	array            unsafe.Pointer
	arraySize        int
	arrayElementSize uintptr
	current          int
}

var _ c.Iterator[any] = (*Iter[any])(nil)

func (i *Iter[T]) HasNext() bool {
	if HasNextByLen(i.arraySize, i.current) {
		i.current++
		return true
	}
	return false
}

func (i *Iter[T]) HasPrev() bool {
	if HasPrevByLen(i.arraySize, i.current) {
		i.current--
		return true
	}
	return false
}

func (i *Iter[T]) Get() T {
	return GetArrayElem[T](i.array, i.current, i.arrayElementSize)
}

type RevertIter[T any] struct {
	Iter[T]
}

var _ c.Iterator[any] = (*RevertIter[any])(nil)

func (i *RevertIter[T]) HasNext() bool {
	if HasPrevByLen(i.arraySize, i.current) {
		i.current--
		return true
	}
	return false
}

func HasNext[T any](elements []T, current int) bool {
	return HasNextByLen(GetLen(elements), current)
}

func HasPrevByLen(size, current int) bool {
	return HasNextByRange(1, size, current)
}

func HasNextByLen(size, current int) bool {
	return HasNextByRange(NoStarted, size-2, current)
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

func GetArrayElementSize[T any](_ []T) uintptr {
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
