package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// IterNoStarted is the head Iterator position
const IterNoStarted = -1

// NewHead instantiates Iter based on elements slice
func NewHead[TS ~[]T, T any](elements TS) Iter[T] {
	var (
		header = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array  = unsafe.Pointer(header.Data)
		size   = header.Len
	)
	return Iter[T]{
		array:   array,
		size:    size,
		current: IterNoStarted,
	}
}

// NewTail instantiates Iter based on elements slice for reverse iterating
func NewTail[T any](elements []T) Iter[T] {
	var (
		header = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array  = unsafe.Pointer(header.Data)
		size   = header.Len
	)
	return Iter[T]{
		array:   array,
		size:    size,
		current: size,
	}
}

// Iter is the Iterator implementation.
type Iter[T any] struct {
	array         unsafe.Pointer
	size, current int
}

var (
	_ c.Iterator[any]     = (*Iter[any])(nil)
	_ c.PrevIterator[any] = (*Iter[any])(nil)
)

func (i *Iter[T]) maxHasNext() int {
	return i.size - 2
}

// All is used to iterate through the iterator using `for ... range`.
func (i *Iter[T]) All(consumer func(element T) bool) {
	loop.All(i.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning Break
func (i *Iter[T]) For(consumer func(element T) error) error {
	return loop.For(i.Next, consumer)
}

// ForEach takes all elements retrieved by the iterator.
func (i *Iter[T]) ForEach(consumer func(element T)) {
	loop.ForEach(i.Next, consumer)
}

// HasNext checks the next element existing
func (i *Iter[T]) HasNext() bool {
	if i == nil {
		return false
	}
	return CanIterateByRange(IterNoStarted, i.maxHasNext(), i.current)
}

// HasPrev checks the previous element existing
func (i *Iter[T]) HasPrev() bool {
	if i == nil {
		return false
	}
	return CanIterateByRange(1, i.size, i.current)
}

// GetNext returns the next element
func (i *Iter[T]) GetNext() T {
	t, _ := i.Next()
	return t
}

// GetPrev returns the previous element
func (i *Iter[T]) GetPrev() T {
	t, _ := i.Prev()
	return t
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *Iter[T]) Next() (v T, ok bool) {
	if !(i == nil || i.array == nil) {
		if current := i.current; CanIterateByRange(IterNoStarted, i.maxHasNext(), current) {
			current++
			i.current = current
			return *(*T)(notsafe.GetArrayElemRef(i.array, current, unsafe.Sizeof(v))), true
		}
	}
	return v, ok
}

// Prev returns the previos element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *Iter[T]) Prev() (v T, ok bool) {
	if !(i == nil || i.array == nil) {
		if current := i.current; CanIterateByRange(1, i.size, current) {
			current--
			i.current = current
			return *(*T)(notsafe.GetArrayElemRef(i.array, current, unsafe.Sizeof(v))), true
		}
	}
	return v, ok
}

// Get returns the current element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *Iter[T]) Get() (v T, ok bool) {
	if !(i == nil || i.array == nil) {
		current := i.current
		if IsValidIndex(i.size, current) {
			return *(*T)(notsafe.GetArrayElemRef(i.array, current, unsafe.Sizeof(v))), true
		}
	}
	return v, ok
}

// Size returns the iterator capacity
func (i *Iter[T]) Size() int {
	if i == nil {
		return 0
	}
	return i.size
}

// Crank rertieves a next element, returns the iterator, element and successfully flag.
func (i *Iter[T]) Crank() (it *Iter[T], t T, ok bool) {
	if i != nil {
		t, ok = i.Next()
	}
	return i, t, ok
}

// CrankPrev rertieves a prev element, returns the iterator, element and successfully flag.
func (i *Iter[T]) CrankPrev() (it *Iter[T], t T, ok bool) {
	if i != nil {
		t, ok = i.Prev()
	}
	return i, t, ok
}

// HasNext checks if an iterator can go forward
func HasNext[T any](elements []T, current int) bool {
	return HasNextBySize(notsafe.GetLen(elements), current)
}

// HasPrev checks if an iterator can go backwards
func HasPrev[T any](elements []T, current int) bool {
	return HasPrevBySize(notsafe.GetLen(elements), current)
}

// HasNextBySize checks if an iterator can go forward
func HasNextBySize(size int, current int) bool {
	return CanIterateByRange(IterNoStarted, size-2, current)
}

// HasPrevBySize checks if an iterator can go backwards
func HasPrevBySize(size, current int) bool {
	return CanIterateByRange(1, size, current)
}

// CanIterateByRange checks if an iterator can go further or stop
func CanIterateByRange(first, last, current int) bool {
	return current >= first && current <= last
}

// IsValidIndex checks if index is out of range
func IsValidIndex(size, index int) bool {
	return index > -1 && index < size
}

// Get safely returns an element of the 'elements' slice by the 'current' index or return zero value of T if the index is more than size-1 or less 0
func Get[TS ~[]T, T any](elements TS, current int) T {
	v, _ := Gett(elements, current)
	return v
}

// Gett safely returns an element of the 'elements' slice by the 'current' index or return zero value of T if the index is more than size-1 or less 0
// ok == true if success
func Gett[TS ~[]T, T any](elements TS, current int) (element T, ok bool) {
	if !(current < 0 || current >= len(elements)) {
		element, ok = (elements)[current], true
	}
	return element, ok
}
