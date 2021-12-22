package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

type ConvertIter[From, To any] struct {
	iter    Iterator[From]
	by      conv.Converter[From, To]
	filters []check.Predicate[From]
	current To
}

var _ Iterator[interface{}] = (*ConvertIter[interface{}, interface{}])(nil)

func (s *ConvertIter[From, To]) Next() (To, bool) {
	v, ok := filter(s.iter, s.filters)
	var r To
	if !ok {
		return r, false
	}
	return s.by(v), true
}

func (s *ConvertIter[From, To]) HasNext() bool {
	v, ok := s.Next()
	s.current = v
	return ok
}

func (s *ConvertIter[From, To]) Get() To {
	return s.current
}