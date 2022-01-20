package it

import (
	"github.com/m4gshm/gollections/typ"
)

type FlattenFit[From, To any] struct {
	Iter  typ.Iterator[From]
	Flatt typ.Flatter[From, To]
	Fit   typ.Predicate[From]

	elementsTo []To
	indTo      int
	c          To
	err        error
}

var _ typ.Iterator[any] = (*FlattenFit[any, any])(nil)

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

	iter := s.Iter
	for {
		if !iter.HasNext() {
			s.err = Exhausted
			return false
		} else if v, err := iter.Get(); err != nil {
			s.err = err
			return true
		} else if s.Fit(v) {
			if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
				s.c = elementsTo[0]
				s.elementsTo = elementsTo
				s.indTo = 1
				return true
			}
		}
	}
}

func (s *FlattenFit[From, To]) Get() (To, error) {
	return s.c, s.err
}

func (s *FlattenFit[From, To]) Next() To {
	return Next[To](s)
}

type Flatten[From, To any] struct {
	Iter  typ.Iterator[From]
	Flatt typ.Flatter[From, To]

	elementsTo []To
	indTo      int
	c          To
	err        error
}

var _ typ.Iterator[any] = (*Flatten[any, any])(nil)

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

	iter := s.Iter
	for {
		if ok := iter.HasNext(); !ok {
			s.err = Exhausted
			return false
		} else if v, err := iter.Get(); err != nil {
			s.err = err
			return true
		} else if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
			s.c = elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return true
		}
	}
}

func (s *Flatten[From, To]) Get() (To, error) {
	return s.c, s.err
}

func (s *Flatten[From, To]) Next() To {
	return Next[To](s)
}
