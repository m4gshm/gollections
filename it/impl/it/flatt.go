package it

import "github.com/m4gshm/gollections/c"

type FlattenFit[From, To any] struct {
	Iter  c.Iterator[From]
	Flatt c.Flatter[From, To]
	Fit   c.Predicate[From]

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

	iter := s.Iter
	for v, ok := iter.Next(); ok && s.Fit(v); v, ok = iter.Next() {
		if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
			c := elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return c, true
		}
	}
	var no To
	return no, false
}

type Flatten[From, To any] struct {
	Iter  c.Iterator[From]
	Flatt c.Flatter[From, To]

	elementsTo []To
	indTo      int
}

var _ c.Iterator[any] = (*Flatten[any, any])(nil)

func (s *Flatten[From, To]) Next() (To, bool) {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.indTo = indTo + 1
			return c, true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	for {
		if v, ok := s.Iter.Next(); !ok {
			var no To
			return no, false
		} else if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
			c := elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return c, true
		}
	}
}
