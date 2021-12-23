package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/slice"
)

type FlattFilterIter[From, To any] struct {
	iter    Iterator[From]
	by      slice.Flatter[From, To]
	filters []check.Predicate[From]
	iterTo  []To
	indTo   int
	c       To
}

var _ Iterator[interface{}] = (*FlattFilterIter[interface{}, interface{}])(nil)

func (s *FlattFilterIter[From, To]) HasNext() bool {
	if iterTo := s.iterTo; len(iterTo) > 0 {
		if indTo := s.indTo; indTo < len(iterTo) {
			c := iterTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.iterTo = nil
		}
	}

	iter := s.iter
	for {
		if ok := iter.HasNext(); !ok {
			return false
		}
		v := iter.Get()
		if ok := check.IsFit(v, s.filters...); ok {
			if iterTo := s.by(v); len(iterTo) > 0 {
				s.c = iterTo[0]
				s.iterTo = iterTo
				s.indTo = 1
				return true
			}
		}
	}
}

func (s *FlattFilterIter[From, To]) Get() To {
	return s.c
}


type FlattIter[From, To any] struct {
	iter   Iterator[From]
	by     slice.Flatter[From, To]
	iterTo []To
	indTo  int
	c      To
}

var _ Iterator[interface{}] = (*FlattIter[interface{}, interface{}])(nil)

func (s *FlattIter[From, To]) HasNext() bool {
	if iterTo := s.iterTo; len(iterTo) > 0 {
		if indTo := s.indTo; indTo < len(iterTo) {
			c := iterTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.iterTo = nil
		}
	}

	iter := s.iter
	for {
		if ok := iter.HasNext(); !ok {
			return false
		}
		v := iter.Get()
		if iterTo := s.by(v); len(iterTo) > 0 {
			s.c = iterTo[0]
			s.iterTo = iterTo
			s.indTo = 1
			return true
		}
	}
}

func (s *FlattIter[From, To]) Get() To {
	return s.c
}