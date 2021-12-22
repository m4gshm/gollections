package iterator

import (
	"github.com/m4gshm/container/check"
)

type FlattIter[From, To any] struct {
	iter    Iterator[From]
	by      Flatter[From, To, Iterator[To]]
	filters []check.Predicate[From]
	current Iterator[To]
	c       To
}

var _ Iterator[interface{}] = (*FlattIter[interface{}, interface{}])(nil)

func (s *FlattIter[From, To]) Next() (To, bool) {
	next := func() (To, bool) {
		n, ok := s.current.Next()

		if !ok {
			s.current = nil
		}
		return n, ok
	}

	if s.current != nil {
		if c, ok := next(); ok {
			return c, true
		}
	}

	var c To
	v, ok := s.iter.Next()
	if !ok {
		return c, false
	}

	ok = check.IsFit(v, s.filters...)
	if !ok {
		return c, false
	}

	current := s.by(v)
	if current != nil {
		s.current = current
		return next()
	} else {
		return c, false
	}
}

func (s *FlattIter[From, To]) HasNext() bool {
	v, ok := s.Next()
	s.c = v
	return ok
}

func (s *FlattIter[From, To]) Get() To {
	return s.c
}
