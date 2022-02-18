package slice

import (
	"github.com/m4gshm/gollections/c"
)

type ConvertFit[From, To any] struct {
	Elements []From
	By       c.Converter[From, To]
	Fit      c.Predicate[From]
	i        int
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

func (s *ConvertFit[From, To]) GetNext() (To, bool) {
	if v, ok := nextArrayElem(s.Elements, s.Fit, &s.i); ok {
		return s.By(v), true
	}
	var no To
	return no, false
}

type Convert[From, To any] struct {
	Elements []From
	By       c.Converter[From, To]
	i        int
	err      error
}

var _ c.Iterator[any] = (*Convert[any, any])(nil)

func (s *Convert[From, To]) GetNext() (To, bool) {
	e := s.Elements
	l := len(s.Elements)
	i := s.i
	if i < l {
		v := e[i]
		s.i = i + 1
		return s.By(v), true
	}
	var no To
	return no, false
}


func nextArrayElem[T any](elements []T, filter c.Predicate[T], indexHolder *int) (T, bool) {
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
