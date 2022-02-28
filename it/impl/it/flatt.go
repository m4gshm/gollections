package it

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

type FlattenFit[From, To any, IT c.Iterator[From]] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	iter          IT
	flatt         c.Flatter[From, To]
	fit           c.Predicate[From]
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
		} else if s.fit(v) {
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

func (s *FlattenFit[From, To, IT]) Cap() int {
	return s.iter.Cap()
}

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

func (s *Flatten[From, To, IT]) Cap() int {
	return s.iter.Cap()
}
