package iter

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/typ"
)

type ConvertFit[From, To any] struct {
	Iter    typ.Iterator[From]
	By      conv.Converter[From, To]
	Fit     check.Predicate[From]
	current To
	err     error
}

var _ typ.Iterator[interface{}] = (*ConvertFit[interface{}, interface{}])(nil)

func (s *ConvertFit[From, To]) HasNext() bool {
	if v, ok := nextFiltered(s.Iter, s.Fit); ok {
		s.current = s.By(v)
		return true
	}
	return false
}

func (s *ConvertFit[From, To]) Get() To {
	return s.current
}

type Convert[From, To any] struct {
	Iter typ.Iterator[From]
	By   conv.Converter[From, To]
}

var _ typ.Iterator[interface{}] = (*Convert[interface{}, interface{}])(nil)

func (s *Convert[From, To]) HasNext() bool {
	return s.Iter.HasNext()
}

func (s *Convert[From, To]) Get() To {
	return s.By(s.Iter.Get())
}
