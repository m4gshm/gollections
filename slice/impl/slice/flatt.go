package slice

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/typ"
)


type FlattenFit[From, To any] struct {
	Elements   []From
	Flatt      conv.Flatter[From, To]
	Fit        check.Predicate[From]

	indFrom    int
	elementsTo []To
	indTo      int
	c          To
}

var _ typ.Iterator[interface{}] = (*FlattenFit[interface{}, interface{}])(nil)

func (s *FlattenFit[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.c = c
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
				s.c = elementsTo[0]
				s.elementsTo = elementsTo
				s.indTo = 1
				return true
			}
		}
	}
	return false
}

func (s *FlattenFit[From, To]) Get() To {
	return s.c
}

type Flatten[From, To any] struct {
	Elements   []From
	Flatt      conv.Flatter[From, To]

	indFrom    int
	elementsTo []To
	indTo      int
	c          To
}

var _ typ.Iterator[interface{}] = (*Flatten[interface{}, interface{}])(nil)

func (s *Flatten[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.c = c
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
			s.c = elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return true
		}
	}
	return false
}

func (s *Flatten[From, To]) Get() To {
	return s.c
}
