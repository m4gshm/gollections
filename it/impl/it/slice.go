package it

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

//NoStarted is the head Iterator position.
const NoStarted = -1

func NewHead[T any, TS ~[]T](elements TS) Iter[T] {
	return NewHeadS(elements, notsafe.GetTypeSize[T]())
}

func NewHeadS[T any, TS ~[]T](elements TS, elementSize uintptr) Iter[T] {
	var (
		header = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array  = unsafe.Pointer(header.Data)
		size   = header.Len
	)
	return Iter[T]{
		array:       array,
		elementSize: elementSize,
		size:        size,
		maxHasNext:  size - 2,
		current:     NoStarted,
	}
}

func NewTail[T any](elements []T) Iter[T] {
	return NewTailS(elements, notsafe.GetTypeSize[T]())
}

func NewTailS[T any](elements []T, elementSize uintptr) Iter[T] {
	var (
		header = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array  = unsafe.Pointer(header.Data)
		size   = header.Len
	)
	return Iter[T]{
		array:       array,
		elementSize: elementSize,
		size:        size,
		maxHasNext:  size - 2,
		current:     size,
	}
}

//Iter is the Iterator implementation.
type Iter[T any] struct {
	array                     unsafe.Pointer
	elementSize               uintptr
	size, maxHasNext, current int
}

var (
	_ c.Iterator[any]     = (*Iter[any])(nil)
	_ c.PrevIterator[any] = (*Iter[any])(nil)
)

func (i *Iter[T]) HasNext() bool {
	return CanIterateByRange(NoStarted, i.maxHasNext, i.current)
}

func (i *Iter[T]) HasPrev() bool {
	return CanIterateByRange(1, i.size, i.current)
}

func (i *Iter[T]) GetNext() T {
	t, _ := i.Next()
	return t
}

func (i *Iter[T]) Next() (T, bool) {
	current := i.current
	if CanIterateByRange(NoStarted, i.maxHasNext, current) {
		current++
		i.current = current
		return *(*T)(notsafe.GetArrayElemRef(i.array, current, i.elementSize)), true
	}
	var no T
	return no, false
}

func (i *Iter[T]) GetPrev() T {
	t, _ := i.Prev()
	return t
}

func (i *Iter[T]) Prev() (T, bool) {
	current := i.current
	if CanIterateByRange(1, i.size, current) {
		current--
		i.current = current
		return *(*T)(notsafe.GetArrayElemRef(i.array, current, i.elementSize)), true
	}
	var no T
	return no, false
}

func (i *Iter[T]) Get() (T, bool) {
	current := i.current
	if IsValidIndex(i.size, current) {
		return *(*T)(notsafe.GetArrayElemRef(i.array, current, i.elementSize)), true
	}
	var no T
	return no, false
}

func (i *Iter[T]) Cap() int {
	return i.size
}

//go:nosplit
func (i Iter[T]) R() *Iter[T] {
	return notsafe.Noescape(&i)
}

//HasNext checks the next element in an iterator by indexs of a current element and slice length.
func HasNext[T any](elements []T, current int) bool {
	return HasNextBySize(notsafe.GetLen(elements), current)
}

//HasPrev checks the previos element in an iterator by indexs of a current element and slice length.
func HasPrev[T any](elements []T, current int) bool {
	return HasPrevBySize(notsafe.GetLen(elements), current)
}

//HasNextBySize checks the next element in an iterator by indexs of a current element and slice length.
func HasNextBySize(size int, current int) bool {
	return CanIterateByRange(NoStarted, size-2, current)
}

//HasPrevBySize checks the previos element in an iterator by indexs of a current element and slice length.
func HasPrevBySize(size, current int) bool {
	return CanIterateByRange(1, size, current)
}

//CanIterateByRange checks the next element in an iterator by indexes of the first, the last and a current elements of an underlying slice.
func CanIterateByRange(first, last, current int) bool {
	return current >= first && current <= last
}

//IsValidIndex checks if index is out of range
func IsValidIndex(size, index int) bool {
	return index > -1 && index < size
}

func IsValidIndex2(size, index int) bool {
	return !((index^size == 0) || index < 0)
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
