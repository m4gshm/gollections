package slice

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/typ"
)

type FlattenFit[From, To any] struct {
	Elements []From
	Flatt    typ.Flatter[From, To]
	Fit      typ.Predicate[From]

	indFrom int

	elementsTo []To
	indTo      int
	current    To
	err        error
}

var _ typ.Iterator[any] = (*FlattenFit[any, any])(nil)

func (s *FlattenFit[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.current = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.elementsTo = nil
		}
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
	s.err = it.Exhausted
	return false
}

func (s *FlattenFit[From, To]) Get() (To, error) {
	return s.current, s.err
}

func (s *FlattenFit[From, To]) Next() To {
	return it.Next[To](s)
}

type Flatten[From, To any] struct {
	Elements []From
	Flatt    typ.Flatter[From, To]

	indFrom    int
	elementsTo []To
	indTo      int
	current    To
	err        error
}

var _ typ.Iterator[any] = (*Flatten[any, any])(nil)

func (s *Flatten[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.current = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.elementsTo = nil
		}
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
	s.err = it.Exhausted
	return false
}

func (s *Flatten[From, To]) Get() (To, error) {
	return s.current, s.err
}

func (s *Flatten[From, To]) Next() To {
	return it.Next[To](s)
}
