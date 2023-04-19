package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

// NoStarted is the head Iterator position.
const NoStarted = -1

// New instantiates Iter based on elements Iter and returs its reference
func New[TS ~[]T, T any](elements TS) ArrayIter[T] {
	return NewHeadS(elements, notsafe.GetTypeSize[T]())
}

// NewHead instantiates Iter based on elements slice
func NewHead[TS ~[]T, T any](elements TS) ArrayIter[T] {
	return NewHeadS(elements, notsafe.GetTypeSize[T]())
}

// NewHeadS instantiates Iter based on elements slice with predefined element size
func NewHeadS[TS ~[]T, T any](elements TS, elementSize uintptr) ArrayIter[T] {
	var (
		header = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array  = unsafe.Pointer(header.Data)
		size   = header.Len
	)
	return ArrayIter[T]{
		array:       array,
		elementSize: elementSize,
		size:        size,
		maxHasNext:  size - 2,
		current:     NoStarted,
	}
}

// NewTail instantiates Iter based on elements slice for reverse iterating
func NewTail[T any](elements []T) ArrayIter[T] {
	return NewTailS(elements, notsafe.GetTypeSize[T]())
}

// NewTailS instantiates Iter based on elements slice with predefined element size for reverse iterating
func NewTailS[T any](elements []T, elementSize uintptr) ArrayIter[T] {
	var (
		header = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array  = unsafe.Pointer(header.Data)
		size   = header.Len
	)
	return ArrayIter[T]{
		array:       array,
		elementSize: elementSize,
		size:        size,
		maxHasNext:  size - 2,
		current:     size,
	}
}

// ArrayIter is the Iterator implementation.
type ArrayIter[T any] struct {
	array                     unsafe.Pointer
	elementSize               uintptr
	size, maxHasNext, current int
}

var (
	_ c.Iterator[any]     = (*ArrayIter[any])(nil)
	_ c.PrevIterator[any] = (*ArrayIter[any])(nil)
)

func (i *ArrayIter[T]) HasNext() bool {
	if i == nil {
		return false
	}
	return CanIterateByRange(NoStarted, i.maxHasNext, i.current)
}

func (i *ArrayIter[T]) HasPrev() bool {
	if i == nil {
		return false
	}
	return CanIterateByRange(1, i.size, i.current)
}

func (i *ArrayIter[T]) GetNext() T {
	t, _ := i.Next()
	return t
}

func (i *ArrayIter[T]) GetPrev() T {
	t, _ := i.Prev()
	return t
}

func (i *ArrayIter[T]) Next() (v T, ok bool) {
	if !(i == nil || i.array == nil) {
		if current := i.current; CanIterateByRange(NoStarted, i.maxHasNext, current) {
			current++
			i.current = current
			return *(*T)(notsafe.GetArrayElemRef(i.array, current, i.elementSize)), true
		}
	}
	return v, ok
}

func (i *ArrayIter[T]) Prev() (v T, ok bool) {
	if !(i == nil || i.array == nil) {
		current := i.current
		if CanIterateByRange(1, i.size, current) {
			current--
			i.current = current
			return *(*T)(notsafe.GetArrayElemRef(i.array, current, i.elementSize)), true
		}
	}
	return v, ok
}

func (i *ArrayIter[T]) Get() (v T, ok bool) {
	if !(i == nil || i.array == nil) {
		current := i.current
		if IsValidIndex(i.size, current) {
			return *(*T)(notsafe.GetArrayElemRef(i.array, current, i.elementSize)), true
		}
	}
	return v, ok
}

func (i *ArrayIter[T]) Cap() int {
	if i == nil {
		return 0
	}
	return i.size
}

// HasNext checks the next element in an iterator by indexs of a current element and slice length.
func HasNext[T any](elements []T, current int) bool {
	return HasNextBySize(notsafe.GetLen(elements), current)
}

// HasPrev checks the previos element in an iterator by indexs of a current element and slice length.
func HasPrev[T any](elements []T, current int) bool {
	return HasPrevBySize(notsafe.GetLen(elements), current)
}

// HasNextBySize checks the next element in an iterator by indexs of a current element and slice length.
func HasNextBySize(size int, current int) bool {
	return CanIterateByRange(NoStarted, size-2, current)
}

// HasPrevBySize checks the previos element in an iterator by indexs of a current element and slice length.
func HasPrevBySize(size, current int) bool {
	return CanIterateByRange(1, size, current)
}

// CanIterateByRange checks the next element in an iterator by indexes of the first, the last and a current elements of an underlying slice.
func CanIterateByRange(first, last, current int) bool {
	return current >= first && current <= last
}

// IsValidIndex checks if index is out of range
func IsValidIndex(size, index int) bool {
	return index > -1 && index < size
}

// IsValidIndex2 checks if index is out of range
func IsValidIndex2(size, index int) bool {
	return !((index^size == 0) || index < 0)
}

// Get safely returns an element of a slice by an index or zero value of T if the index is out of range.
func Get[TS ~[]T, T any](elements TS, current int) T {
	v, _ := Gett(elements, current)
	return v
}

// Gett safely returns an element of a slice adn true by an index or zero value of T and false if the index is out of range.
func Gett[TS ~[]T, T any](elements TS, current int) (T, bool) {
	if current >= len(elements) {
		var no T
		return no, false
	} else if current == NoStarted {
		var no T
		return no, false
	}
	return (elements)[current], true
}
