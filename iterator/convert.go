package iterator

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

type ConvertFilterIter[From, To any] struct {
	iter    Iterator[From]
	by      conv.Converter[From, To]
	filters []check.Predicate[From]
	current To
}

var _ Iterator[interface{}] = (*ConvertFilterIter[interface{}, interface{}])(nil)

func (s *ConvertFilterIter[From, To]) HasNext() bool {
	if v, ok := filterNext(s.iter, s.filters); ok {
		s.current = s.by(v)
		return true
	}
	return false
}

func (s *ConvertFilterIter[From, To]) Get() To {
	return s.current
}

type ConvertIter[From, To any] struct {
	iter    Iterator[From]
	by      conv.Converter[From, To]
}

var _ Iterator[interface{}] = (*ConvertIter[interface{}, interface{}])(nil)

func (s *ConvertIter[From, To]) HasNext() bool {
	return s.iter.HasNext()
}

func (s *ConvertIter[From, To]) Get() To { 
	return s.by(s.iter.Get())
}