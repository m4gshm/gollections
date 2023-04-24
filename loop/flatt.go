package loop

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

// FlattenFitIter is the Iterator wrapper that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlattenFitIter[From, To any] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	next          func() (From, bool)
	flatt         func(From) []To
	filterFrom    func(From) bool
	filterTo      func(To) bool
}

var _ c.Iterator[any] = (*FlattenFitIter[any, any])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *FlattenFitIter[From, To]) Next() (t To, ok bool) {
	if s == nil {
		return t, ok
	}

	next := s.next
	if next == nil {
		return t, false
	}

	for {
		if sizeTo := s.sizeTo; sizeTo > 0 {
			if indTo := s.indTo; indTo < sizeTo {
				s.indTo++
				t = *(*To)(notsafe.GetArrayElemRef(s.arrayTo, indTo, s.elemSizeTo))
				if ok := s.filterTo(t); ok {
					return t, true
				}
			}
			s.indTo = 0
			s.arrayTo = nil
			s.sizeTo = 0
		}

		if v, ok := next(); !ok {
			return t, false
		} else if s.filterFrom(v) {
			if elementsTo := s.flatt(v); len(elementsTo) > 0 {
				s.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				s.arrayTo = unsafe.Pointer(header.Data)
				s.sizeTo = header.Len
				t = *(*To)(notsafe.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo))
				if ok := s.filterTo(t); ok {
					return t, true
				}
			}
		}
	}
}

// FlatIter is the Iterator wrapper that converts an element to a slice and iterates over the elements of that slice.
// For example, FlatIter can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type FlatIter[From, To any] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	next          func() (From, bool)
	flatt         func(From) []To
}

var _ c.Iterator[any] = (*FlatIter[any, any])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *FlatIter[From, To]) Next() (To, bool) {
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
		if v, ok := s.next(); !ok {
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
