package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

// ConvertFit is the array based Iterator thath provides converting of elements by a Converter with addition filtering of the elements by a Predicate.
type ConvertFit[From, To any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	by       func(From) To
	filter   func(From) bool
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

func (s *ConvertFit[From, To]) Next() (To, bool) {
	if v, ok := nextFiltered(s.array, s.size, s.elemSize, s.filter, &s.i); ok {
		return s.by(v), true
	}
	var no To
	return no, false
}

func (s *ConvertFit[From, To]) Cap() int {
	return s.size
}

// Converter is the array based Iterator thath provides converting of elements by a Converter.
type Converter[From, To any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	by       func(From) To
}

var _ c.Iterator[any] = (*Converter[any, any])(nil)

func (s *Converter[From, To]) Next() (To, bool) {
	if s.i < s.size {
		v := *(*From)(notsafe.GetArrayElemRef(s.array, s.i, s.elemSize))
		s.i++
		return s.by(v), true
	}
	var no To
	return no, false
}

func (s *Converter[From, To]) Cap() int {
	return s.size
}

func nextFiltered[T any](array unsafe.Pointer, size int, elemSize uintptr, filter func(T) bool, index *int) (T, bool) {
	for i := *index; i < size; i++ {
		if v := *(*T)(notsafe.GetArrayElemRef(array, i, elemSize)); filter(v) {
			*index = i + 1
			return v, true
		}
	}
	var v T
	return v, false
}
