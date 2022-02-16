package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

//NoStarted is the head Iterator position.
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

//Iter is the Iterator implementation.
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

//HasNext checks the next element in an iterator by indexs of a current element and slice length.
func HasNext[T any](elements []T, current int) bool {
	return HasNextBySize(GetLen(elements), current)
}

//HasPrev checks the previos element in an iterator by indexs of a current element and slice length.
func HasPrev[T any](elements []T, current int) bool {
	return HasPrevBySize(GetLen(elements), current)
}

//HasNextBySize checks the next element in an iterator by indexs of a current element and slice length.
func HasNextBySize(size int, current int) bool {
	return HasNextByRange(NoStarted, size-2, current)
}

//HasPrevBySize checks the previos element in an iterator by indexs of a current element and slice length.
func HasPrevBySize(size, current int) bool {
	return HasNextByRange(1, size, current)
}

//HasNextByRange checks the next element in an iterator by indexes of the first, the last and a current elements of an underlying slice.
func HasNextByRange(first, last, current int) bool {
	return current >= first && current <= last
}

//Get safely returns an element of a slice by an index or zero value of T if the index is out of range.
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

//GetArrayPointer retrieves the pointer of a slice underlying array
func GetArrayPointer[T any](elements []T) unsafe.Pointer {
	return unsafe.Pointer(GetSliceHeader(elements).Data)
}

//GetTypeSize retrieves size of a type T
func GetTypeSize[T any]() uintptr {
	var t T
	return unsafe.Sizeof(t)
}

//GetArrayElem returns an element by index from the array referenced by an unsafe pointer
func GetArrayElem[T any](array unsafe.Pointer, index int, elemSize uintptr) T {
	return *(*T)(unsafe.Pointer(uintptr(array) + uintptr(index)*elemSize))
}

//GetLen returns the length of elements
func GetLen[T any](elements []T) int {
	return GetSliceHeader(elements).Len
}

//GetSliceHeader retrieves the SliceHeader of elements
func GetSliceHeader[T any](elements []T) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&elements))
}
