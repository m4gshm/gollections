package it

import "github.com/m4gshm/gollections/c"

type FlattenFit[From, To any] struct {
	Iter  c.Iterator[From]
	Flatt c.Flatter[From, To]
	Fit   c.Predicate[From]

	elementsTo []To
	indTo      int
	c          To
	err        error
}

var _ c.Iterator[any] = (*FlattenFit[any, any])(nil)

func (s *FlattenFit[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	iter := s.Iter
	for {
		if !iter.HasNext() {
			return false
		} else if v := iter.Get(); s.Fit(v) {
			if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
				s.c = elementsTo[0]
				s.elementsTo = elementsTo
				s.indTo = 1
				return true
			}
		}
	}
}

func (s *FlattenFit[From, To]) Get() To {
	return s.c
}

type Flatten[From, To any] struct {
	Iter  c.Iterator[From]
	Flatt c.Flatter[From, To]

	elementsTo []To
	indTo      int
	c          To
}

var _ c.Iterator[any] = (*Flatten[any, any])(nil)

func (s *Flatten[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	iter := s.Iter
	for {
		if ok := iter.HasNext(); !ok {
			return false
		} else if elementsTo := s.Flatt(iter.Get()); len(elementsTo) > 0 {
			s.c = elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return true
		}
	}
}

func (s *Flatten[From, To]) Get() To {
	return s.c
}
