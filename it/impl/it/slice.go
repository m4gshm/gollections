package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

//NoStarted is the head Iterator position.
const NoStarted = -1

func NewHead[T any, TS ~[]T](elements TS) Iter[T] {
	return NewHeadS[T](elements, GetTypeSize[T]())
}

func NewHeadS[T any, TS ~[]T](elements TS, elementSize uintptr) Iter[T] {
	header := GetSliceHeader[T](elements)
	return Iter[T]{
		array:       unsafe.Pointer(header.Data),
		elementSize: elementSize,
		size:        header.Len,
		current:     NoStarted,
	}
}

func NewTail[T any](elements []T) Iter[T] {
	return NewTailS(elements, GetTypeSize[T]())
}

func NewTailS[T any](elements []T, elementSize uintptr) Iter[T] {
	var (
		header = GetSliceHeader(elements)
		size   = header.Len
	)
	return Iter[T]{
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

var (
	_ c.Iterator[any]     = (*Iter[any])(nil)
	_ c.PrevIterator[any] = (*Iter[any])(nil)
)

func (i *Iter[T]) HasNext() bool {
	return HasNextBySize(i.size, i.current)
}

func (i *Iter[T]) HasPrev() bool {
	return HasPrevBySize(i.size, i.current)
}

func (i *Iter[T]) Next() T {
	t, _ := i.GetNext()
	return t
}

func (i *Iter[T]) GetNext() (T, bool) {
	if i.HasNext() {
		i.current++
		return GetArrayElem[T](i.array, i.current, i.elementSize), true
	}
	var no T
	return no, false
}

func (i *Iter[T]) Prev() T {
	t, _ := i.GetPrev()
	return t
}

func (i *Iter[T]) GetPrev() (T, bool) {
	if i.HasPrev() {
		i.current--
		return GetArrayElem[T](i.array, i.current, i.elementSize), true
	}
	var no T
	return no, false
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
	v, _ := Gett(elements, current)
	return v
}

//Gett safely returns an element of a slice adn true by an index or zero value of T and false if the index is out of range.
func Gett[T any](elements []T, current int) (T, bool) {
	if current >= len(elements) {
		var no T
		return no, false
	} else if current == NoStarted {
		var no T
		return no, false
	}
	return (elements)[current], true
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
func GetSliceHeader[T any, TS ~[]T](elements TS) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&elements))
}
