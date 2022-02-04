package it

import (
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

const NoStarted = -1

func New[T any](elements []T) *Iter[T] {
	return NewPP[T](GetArrayPointer(&elements), GetLen(&elements), GetArrayElementSize(&elements))
}

func NewPP[T any](array unsafe.Pointer, arraySize int, arrayElementSize uintptr) *Iter[T] {
	return &Iter[T]{array: array, arraySize: arraySize, arrayElementSize: arrayElementSize, current: NoStarted}
}

type Iter[T any] struct {
	array            unsafe.Pointer
	arraySize        int
	arrayElementSize uintptr
	current          int
}

var _ c.Iterator[any] = (*Iter[any])(nil)
var _ c.Resetable = (*Iter[any])(nil)

func (s *Iter[T]) HasNext() bool {
	return HasNextByLen(s.arraySize, &s.current)
}

func (s *Iter[T]) Next() T {
	return *GetArrayElem[T](s.array, s.current, s.arrayElementSize)
}

func (s *Iter[T]) Position() int {
	return s.current
}

func (s *Iter[T]) SetPosition(pos int) {
	s.current = pos
}

func (s *Iter[T]) Reset() {
	s.SetPosition(NoStarted)
}

type sliceType struct {
	rtype
	elem *rtype // slice element type
}

type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type textOff int32 // offset from top of text section

type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte   // garbage collection data
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}

func HasNext[T any](elements *[]T, current *int) bool {
	return HasNextByLen(GetLen(elements), current)
}

func HasNextByLen(size int, index *int) bool {
	if size == 0 {
		return false
	}
	c := *index
	if c == NoStarted || c < (size-1) {
		*index++
		return true
	}
	return false
}

func Get[T any](elements *[]T, current int) T {
	if current >= len(*elements) {
		var no T
		return no
	} else if current == NoStarted {
		var no T
		return no
	}
	return (*elements)[current]
}

func GetArrayPointer[T any](elements *[]T) unsafe.Pointer {
	return unsafe.Pointer(GetSliceHeader(elements).Data)
}

func GetArrayElementSize[T any](_ *[]T) uintptr {
	var t T
	return unsafe.Sizeof(t)
}

func GetArrayElem[T any](array unsafe.Pointer, index int, elemSize uintptr) *T {
	return (*T)(unsafe.Pointer(uintptr(array) + uintptr(index)*elemSize))
}

func GetLen[T any](elements *[]T) int {
	return GetSliceHeader(elements).Len
}

func GetSliceHeader[T any](elements *[]T) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(elements))
}
