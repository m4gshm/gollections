package it

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	sunsafe "github.com/m4gshm/gollections/slice/unsafe"
)

type FlattenFit[From, To any, IT c.Iterator[From]] struct {
	iter       IT
	flatt      c.Flatter[From, To]
	fit        c.Predicate[From]
	elementsTo []To
	indTo      int
}

var _ c.Iterator[any] = (*FlattenFit[any, any, c.Iterator[any]])(nil)

func (s *FlattenFit[From, To, IT]) Next() (To, bool) {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.indTo = indTo + 1
			return c, true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	iter := s.iter
	for v, ok := iter.Next(); ok && s.fit(v); v, ok = iter.Next() {
		if elementsTo := s.flatt(v); len(elementsTo) > 0 {
			c := elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return c, true
		}
	}
	var no To
	return no, false
}

func (s *FlattenFit[From, To, IT]) Cap() int {
	return s.iter.Cap()
}

func (s FlattenFit[From, To, IT]) R() *FlattenFit[From, To, IT] {
	return (*FlattenFit[From, To, IT])(noescape(&s))
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
			return *(*To)(sunsafe.GetArrayElemRef(s.arrayTo, indTo, s.elemSizeTo)), true
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
			header := sunsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			s.arrayTo = unsafe.Pointer(header.Data)
			s.sizeTo = header.Len
			return *(*To)(sunsafe.GetArrayElemRef(s.arrayTo, 0, s.elemSizeTo)), true
		}
	}
}

func (s *Flatten[From, To, IT]) Cap() int {
	return s.iter.Cap()
}

func (s Flatten[From, To, IT]) R() *Flatten[From, To, IT] {
	return (*Flatten[From, To, IT])(noescape(&s))
}

//go:nosplit
//go:nocheckptr
func noescape[T any](t *T) unsafe.Pointer {
	x := uintptr(unsafe.Pointer(t))
	return unsafe.Pointer(x ^ 0)
}
