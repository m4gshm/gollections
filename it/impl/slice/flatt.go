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

	elements := s.Elements
	le := len(elements)
	for indFrom := s.indFrom; indFrom < le; indFrom++ {
		s.indFrom = indFrom + 1
		if v := elements[indFrom]; s.Fit(v) {
			if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
				c := elementsTo[0]
				s.elementsTo = elementsTo
				s.indTo = 1
				return c, true
			}
		}
	}
	var no To
	return no, false
}

//Flatten is the Iterator impelementation that converts an element to a slice.
//For example, Flatten can be used to convert a multi-dimensional array to a one-dimensional array ([][]int -> []int).
type Flatten[From, To any] struct {
	Elements []From
	Flatt    c.Flatter[From, To]

	indFrom, indTo *int
	elementsTo     []To
}

var _ c.Iterator[any] = (*Flatten[any, any])(nil)

func (s *Flatten[From, To]) Next() (To, bool) {
	if s.elementsTo != nil && len(s.elementsTo) > 0 {
		if indTo := *s.indTo; indTo < len(s.elementsTo) {
			c := (s.elementsTo)[indTo]
			*s.indTo = indTo + 1
			return c, true
		}
		*s.indTo = 0
		s.elementsTo = nil
	}

	elements := s.Elements
	le := len(elements)
	for indFrom := *s.indFrom; indFrom < le; indFrom++ {
		*s.indFrom = indFrom + 1
		v := elements[indFrom]
		if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
			c := elementsTo[0]
			s.elementsTo = elementsTo
			*s.indTo = 1
			return c, true
		}
	}
	var no To
	return no, false
}
