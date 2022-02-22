package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

type FlattenFit[From, To any] struct {
	elements            []From
	flatt               c.Flatter[From, To]
	fit                 c.Predicate[From]
	elementsTo          []To
	indFrom, indTo, cap int
}

var _ c.Iterator[any] = (*FlattenFit[any, any])(nil)

func (s *FlattenFit[From, To]) Next() (To, bool) {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.indTo = indTo + 1
			return c, true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	elements := s.elements
	le := len(elements)
	for indFrom := s.indFrom; indFrom < le; indFrom++ {
		s.indFrom = indFrom + 1
		if v := elements[indFrom]; s.fit(v) {
			if elementsTo := s.flatt(v); len(elementsTo) > 0 {
				c := elementsTo[0]
				s.elementsTo = elementsTo
				s.cap += len(elementsTo)
				s.indTo = 1
				return c, true
			}
		}
	}
	var no To
	return no, false
}

func (s *FlattenFit[From, To]) Cap() int {
	return s.cap
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
			return *(*To)(it.GetArrayElemRef(s.arrayTo, indTo, s.elemSizeTo)), true
		}
		s.indTo = 0
		s.arrayTo = nil
		s.sizeTo = 0
	}
	for indFrom := s.indFrom; indFrom < s.sizeFrom; indFrom++ {
		if elementsTo := s.flatt(*(*From)(it.GetArrayElemRef(s.arrayFrom, indFrom, s.elemSizeFrom))); len(elementsTo) > 0 {
			s.indFrom = indFrom + 1
			s.indTo = 1
			header := it.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			s.arrayTo = unsafe.Pointer(header.Data)
			s.sizeTo = header.Len
			return  *(*To)(it.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo)), true
		}
	}
	var no To
	return no, false
}

func (s *Flatten[From, To]) Cap() int {
	return s.sizeFrom
}
