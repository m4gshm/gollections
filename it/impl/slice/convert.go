package slice

import (
	"github.com/m4gshm/gollections/c"
)

type ConvertFit[From, To any] struct {
	elements []From
	by       c.Converter[From, To]
	Fit      c.Predicate[From]
	i        int
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

func (s *ConvertFit[From, To]) Next() (To, bool) {
	if v, ok := nextArrayElem(s.elements, s.Fit, &s.i); ok {
		return s.by(v), true
	}
	var no To
	return no, false
}

func (s *ConvertFit[From, To]) Cap() int {
	return len(s.elements)
}

func (s ConvertFit[From, To]) R() *ConvertFit[From, To] {
	return &s
	// return notsafe.Noescape(&s)
}

type Convert[From, To any] struct {
	elements []From
	by       c.Converter[From, To]
	i        int
	err      error
}

var _ c.Iterator[any] = (*Convert[any, any])(nil)

func (s *Convert[From, To]) Next() (To, bool) {
	e := s.elements
	l := len(s.elements)
	i := s.i
	if i < l {
		v := e[i]
		s.i = i + 1
		return s.by(v), true
	}
	var no To
	return no, false
}

func (s *Convert[From, To]) Cap() int {
	return len(s.elements)
}

func (s Convert[From, To]) R() *Convert[From, To] {
	return &s
	// return notsafe.Noescape(&s)
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
