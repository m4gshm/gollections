package slice

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/typ"
)

type ConvertFit[From, To any] struct {
	Elements []From
	By       conv.Converter[From, To]
	Fit      check.Predicate[From]

	i        int
	current  To
}

var _ typ.Iterator[interface{}] = (*ConvertFit[interface{}, interface{}])(nil)

func (s *ConvertFit[From, To]) HasNext() bool {
	if v, ok := nextArrayElem(s.Elements, s.Fit, &s.i); ok {
		s.current = s.By(v)
		return true
	}
	return false
}

func (s *ConvertFit[From, To]) Get() To {
	return s.current
}

type Convert[From, To any] struct {
	Elements []From
	By       conv.Converter[From, To]

	i        int
	current  To

}

var _ typ.Iterator[interface{}] = (*Convert[interface{}, interface{}])(nil)

func (s *Convert[From, To]) HasNext() bool {
	e := s.Elements
	l := len(s.Elements)
	i := s.i
	if i < l {
		v := e[i]
		s.i = i + 1
		s.current = s.By(v)
		return true
	}
	return false
}

func (s *Convert[From, To]) Get() To {
	return s.current
}


func nextArrayElem[T any](elements []T, filter check.Predicate[T], indexHolder *int) (T, bool) {
	l := len(elements)
	for i := *indexHolder; i < l; i++ {
		if v := elements[i]; filter(v) {
			*indexHolder = i + 1
			return v, true
		}
	}
	var v T
	return v, false
}