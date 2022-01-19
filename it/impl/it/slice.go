package it

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/m4gshm/gollections/typ"
)

const NoStarted = -1

var (
	Exhausted        = errors.New("exhausted interator")
	GetBeforeHasNext = errors.New("'Get' called before 'HasNext'")
)

func New[T any](elements []T) *Iter[T] {
	return &Iter[T]{elements: elements, size: len(elements), current: NoStarted}
}

func NewReseteable[T any](elements []T) *Reseteable[T] {
	return &Reseteable[T]{New(elements)}
}

type Iter[T any] struct {
	elements []T
	size     int
	err      error
	current  int
}

var _ typ.Iterator[any] = (*Iter[any])(nil)

func (s *Iter[T]) HasNext() bool {
	size := s.size
	if size == 0 {
		s.err = Exhausted
		return false
	}
	c := s.current
	if c == NoStarted || c < (size-1) {
		s.current++
		return true
	}
	s.err = Exhausted
	return false
}

func (s *Iter[T]) Get() (T, error) {
	r, err := Get(s.current, &s.elements, s.err)
	if err != nil {
		var no T
		return no, err
	}
	return r, nil
}

func (s *Iter[T]) Position() int {
	return s.current
}

func (s *Iter[T]) SetPosition(pos int) {
	s.current = pos
}

type Reseteable[T any] struct {
	*Iter[T]
}

var _ typ.Resetable = (*Reseteable[interface{}])(nil)

func (s *Reseteable[T]) Reset() {
	s.SetPosition(NoStarted)
	s.err = nil
}

func NewP[T any](elements *[]T) *PIter[T] {
	return NewPP[T](GetArrayPointer(elements), GetLen(elements), GetArrayElementSize(elements))
}

func NewPP[T any](array unsafe.Pointer, arraySize int, arrayElementSize uintptr) *PIter[T] {
	return &PIter[T]{array: array, arraySize: arraySize, arrayElementSize: arrayElementSize, current: NoStarted}
}

type PIter[T any] struct {
	array            unsafe.Pointer
	arraySize        int
	arrayElementSize uintptr

	err     error
	current int
}

var _ typ.Iterator[any] = (*PIter[any])(nil)

func (s *PIter[T]) HasNext() bool {
	return HasNextByLen(s.arraySize, &s.current, &s.err)
}

func (s *PIter[T]) Get() (T, error) {
	return *GetArrayElem[T](s.array, s.current, s.arrayElementSize), nil
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

func HasNext[T any](elements *[]T, current *int, err *error) bool {
	return HasNextByLen(GetLen(elements), current, err)
}

func HasNextByLen(size int, index *int, err *error) bool {
	if size == 0 {
		*err = Exhausted
		return false
	}
	c := *index
	if c == NoStarted || c < (size-1) {
		*index++
		return true
	}
	*err = Exhausted
	return false
}

func Get[T any](current int, elements *[]T, err error) (T, error) {
	if err != nil {
		var no T
		return no, err
	} else if current == NoStarted {
		var no T
		return no, GetBeforeHasNext
	}
	return (*elements)[current], nil
}

func GetArrayPointer[T any](elements *[]T) unsafe.Pointer {
	return unsafe.Pointer(GetSliceHeader(elements).Data)
}

func GetArrayElementSize2[T any](elements []T) uintptr {
	a := any(elements)
	typ := *(**sliceType)(unsafe.Pointer(&a))
	return typ.elem.size
}

func GetArrayElementSize[T any](elements *[]T) uintptr {
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
