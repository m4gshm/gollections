package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

type FlattenFit[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo, cap      int
	flatt                    c.Flatter[From, To]
	fit                      c.Predicate[From]
}

var _ c.Iterator[any] = (*FlattenFit[any, any])(nil)

func (s *FlattenFit[From, To]) Next() (To, bool) {
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
		if v := *(*From)(notsafe.GetArrayElemRef(s.arrayFrom, indFrom, s.elemSizeFrom)); s.fit(v) {
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

func (s *FlattenFit[From, To]) Cap() int {
	return s.cap
}

func (s FlattenFit[From, To]) R() *FlattenFit[From, To] {
	return notsafe.Noescape(&s)
}

//Flatten is the Iterator impelementation that converts an element to a slice.
//For example, Flatten can be used to convert a multi-dimensional array to a one-dimensional array ([][]int -> []int).
type Flatten[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo           int
	flatt                    c.Flatter[From, To]
}

var _ c.Iterator[any] = (*Flatten[any, any])(nil)

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

func (s *Flatten[From, To]) Cap() int {
	return s.sizeFrom
}

func (s Flatten[From, To]) R() *Flatten[From, To] {
	return notsafe.Noescape(&s)
}
