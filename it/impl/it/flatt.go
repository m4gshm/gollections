package it

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/predicate"
)

// FlattenFit is the Iterator wrapper that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlattenFit[From, To any, IT c.Iterator[From]] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	iter          IT
	flatt         c.Flatter[From, To]
	filter        predicate.Predicate[From]
}

var _ c.Iterator[any] = (*FlattenFit[any, any, c.Iterator[any]])(nil)

func (s *FlattenFit[From, To, IT]) Next() (To, bool) {
	sizeTo := s.sizeTo
	if sizeTo > 0 {
		if indTo := s.indTo; indTo < sizeTo {
			s.indTo++
			return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, indTo, s.elemSizeTo)), true
		}
		s.indTo = 0
		s.arrayTo = nil
		s.sizeTo = 0
	}

	for {
		if v, ok := s.iter.Next(); !ok {
			var no To
			return no, false
		} else if s.filter(v) {
			if elementsTo := s.flatt(v); len(elementsTo) > 0 {
				s.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				s.arrayTo = unsafe.Pointer(header.Data)
				s.sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo)), true
			}
		}
	}
}

// Flatten is the Iterator wrapper that converts an element to a slice and iterates over the elements of that slice.
// For example, Flatten can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type Flatten[From, To any, IT c.Iterator[From]] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	iter          IT
	flatt         c.Flatter[From, To]
}

var _ c.Iterator[any] = (*Flatten[any, any, c.Iterator[any]])(nil)

func (s *Flatten[From, To, IT]) Next() (To, bool) {
	sizeTo := s.sizeTo
	if sizeTo > 0 {
		if indTo := s.indTo; indTo < sizeTo {
			s.indTo++
			return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, indTo, s.elemSizeTo)), true
		}
		s.indTo = 0
		s.arrayTo = nil
		s.sizeTo = 0
	}

	for {
		if v, ok := s.iter.Next(); !ok {
			var no To
			return no, false
		} else if elementsTo := s.flatt(v); len(elementsTo) > 0 {
			s.indTo = 1
			header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			s.arrayTo = unsafe.Pointer(header.Data)
			s.sizeTo = header.Len
			return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo)), true
		}
	}
}
