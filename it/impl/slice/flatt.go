package slice

import (
	"github.com/m4gshm/gollections/c"
)

type FlattenFit[From, To any] struct {
	Elements []From
	Flatt    c.Flatter[From, To]
	Fit      c.Predicate[From]

	indFrom int

	elementsTo []To
	indTo      int
	current    To
}

var _ c.Iterator[any] = (*FlattenFit[any, any])(nil)

func (s *FlattenFit[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.current = c
			s.indTo = indTo + 1
			return true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	elements := s.Elements
	le := len(elements)
	for indFrom := s.indFrom; indFrom < le; indFrom++ {
		s.indFrom = indFrom + 1
		if v := elements[indFrom]; s.Fit(v) {
			if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
				s.current = elementsTo[0]
				s.elementsTo = elementsTo
				s.indTo = 1
				return true
			}
		}
	}
	return false
}

func (s *FlattenFit[From, To]) Get() To {
	return s.current
}

//Flatten is the Iterator impelementation that converts an element to a slice.
//For example, Flatten can be used to convert a multi-dimensional array to a one-dimensional array ([][]int -> []int).
type Flatten[From, To any] struct {
	Elements []From
	Flatt    c.Flatter[From, To]

	indFrom    int
	elementsTo []To
	indTo      int
	current    To
}

var _ c.Iterator[any] = (*Flatten[any, any])(nil)

func (s *Flatten[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.current = c
			s.indTo = indTo + 1
			return true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	elements := s.Elements
	le := len(elements)
	for indFrom := s.indFrom; indFrom < le; indFrom++ {
		s.indFrom = indFrom + 1
		v := elements[indFrom]
		if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
			s.current = elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return true
		}
	}
	return false
}

func (s *Flatten[From, To]) Get() To {
	return s.current
}
