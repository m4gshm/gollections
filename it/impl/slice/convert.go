package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

type ConvertFit[From, To any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	by       c.Converter[From, To]
	Fit      c.Predicate[From]
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

func (s *ConvertFit[From, To]) Next() (To, bool) {
	if v, ok := nextFiltered(s.array, s.size, s.elemSize, s.Fit, &s.i); ok {
		return s.by(v), true
	}
	var no To
	return no, false
}

func (s *ConvertFit[From, To]) Cap() int {
	return s.size
}

type Convert[From, To any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	by       c.Converter[From, To]
}

var _ c.Iterator[any] = (*Convert[any, any])(nil)

func (s *Convert[From, To]) Next() (To, bool) {
	if s.i < s.size {
		v := *(*From)(notsafe.GetArrayElemRef(s.array, s.i, s.elemSize))
		s.i++
		return s.by(v), true
	}
	var no To
	return no, false
}

func (s *Convert[From, To]) Cap() int {
	return s.size
}

func nextFiltered[T any](array unsafe.Pointer, size int, elemSize uintptr, filter c.Predicate[T], index *int) (T, bool) {
	for i := *index; i < size; i++ {
		if v := *(*T)(notsafe.GetArrayElemRef(array, i, elemSize)); filter(v) {
			*index = i + 1
			return v, true
		}
	}
	var v T
	return v, false
}
