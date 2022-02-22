package it

import "github.com/m4gshm/gollections/c"

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

type Flatten[From, To any, IT c.Iterator[From]] struct {
	iter       IT
	flatt      c.Flatter[From, To]
	elementsTo []To
	indTo      int
}

var _ c.Iterator[any] = (*Flatten[any, any, c.Iterator[any]])(nil)

func (s *Flatten[From, To, IT]) Next() (To, bool) {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.indTo++
			return c, true
		}
		s.indTo = 0
		s.elementsTo = nil
	}

	for {
		if v, ok := s.iter.Next(); !ok {
			var no To
			return no, false
		} else if elementsTo := s.flatt(v); len(elementsTo) > 0 {
			c := elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return c, true
		}
	}
}

func (s *Flatten[From, To, IT]) Cap() int {
	return s.iter.Cap()
}
