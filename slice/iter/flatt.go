package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

// FlattenFitIter is the array based Iterator impelementation that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlattenFitIter[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo, cap      int
	flatt                    func(From) []To
	filter                   func(From) bool
}

var _ c.Iterator[any] = (*FlattenFitIter[any, any])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *FlattenFitIter[From, To]) Next() (To, bool) {
	if sizeTo := s.sizeTo; sizeTo > 0 {
		if indTo := s.indTo; indTo < sizeTo {
			s.indTo++
			return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, indTo, s.elemSizeTo)), true
		}
		s.indTo = 0
		s.arrayTo = nil
		s.sizeTo = 0
	}

	for indFrom := s.indFrom; indFrom < s.sizeFrom; indFrom++ {
		if v := *(*From)(notsafe.GetArrayElemRef(s.arrayFrom, indFrom, s.elemSizeFrom)); s.filter(v) {
			if elementsTo := s.flatt(v); len(elementsTo) > 0 {
				s.indFrom = indFrom + 1
				s.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				s.arrayTo = unsafe.Pointer(header.Data)
				s.sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo)), true
			}
		}
	}
	var no To
	return no, false
}

// Cap returns the iterator capacity
func (s *FlattenFitIter[From, To]) Cap() int {
	return s.cap
}

// Flatten is the array based Iterator impelementation that converts an element to a slice and iterates over the elements of that slice.
// For example, Flatten can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type Flatten[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo           int
	flatt                    func(From) []To
}

var _ c.Iterator[any] = (*Flatten[any, any])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *Flatten[From, To]) Next() (To, bool) {
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
	for indFrom := s.indFrom; indFrom < s.sizeFrom; indFrom++ {
		if elementsTo := s.flatt(*(*From)(notsafe.GetArrayElemRef(s.arrayFrom, indFrom, s.elemSizeFrom))); len(elementsTo) > 0 {
			s.indFrom = indFrom + 1
			s.indTo = 1
			header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			s.arrayTo = unsafe.Pointer(header.Data)
			s.sizeTo = header.Len
			return *(*To)(notsafe.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo)), true
		}
	}
	var no To
	return no, false
}

// Cap returns the iterator capacity
func (s *Flatten[From, To]) Cap() int {
	return s.sizeFrom
}
