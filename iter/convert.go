package iter

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/typ"
)

type ConvertFit[From, To any] struct {
	iter    typ.Iterator[From]
	by      conv.Converter[From, To]
	fit     check.Predicate[From]
	current To
}

var _ typ.Iterator[interface{}] = (*ConvertFit[interface{}, interface{}])(nil)

func (s *ConvertFit[From, To]) HasNext() bool {
	if v, ok := nextFiltered(s.iter, s.fit); ok {
		s.current = s.by(v)
		return true
	}
	return false
}

func (s *ConvertFit[From, To]) Get() To {
	return s.current
}

type Convert[From, To any] struct {
	iter typ.Iterator[From]
	by   conv.Converter[From, To]
}

var _ typ.Iterator[interface{}] = (*Convert[interface{}, interface{}])(nil)

func (s *Convert[From, To]) HasNext() bool {
	return s.iter.HasNext()
}

func (s *Convert[From, To]) Get() To {
	return s.by(s.iter.Get())
}
