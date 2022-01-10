package slice

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/typ"
)

type ConvertFit[From, To any] struct {
	Elements []From
	By       typ.Converter[From, To]
	Fit      typ.Predicate[From]

	i       int
	current To
	err     error
}

var _ typ.Iterator[interface{}] = (*ConvertFit[interface{}, interface{}])(nil)

func (s *ConvertFit[From, To]) HasNext() bool {
	if v, ok := nextArrayElem(s.Elements, s.Fit, &s.i); ok {
		s.current = s.By(v)
		return true
	}
	s.err = it.Exhausted
	return false
}

func (s *ConvertFit[From, To]) Get() To {
	v, err := s.Next()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *ConvertFit[From, To]) Next() (To, error) {
	return s.current, s.err
}

func (s *ConvertFit[From, To]) Err() error {
	return s.err
}

type Convert[From, To any] struct {
	Elements []From
	By       typ.Converter[From, To]

	i       int
	current To
	err     error
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
	s.err = it.Exhausted
	return false
}

func (s *Convert[From, To]) Get() To {
	v, err := s.Next()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *Convert[From, To]) Next() (To, error) {
	return s.current, s.err
}

func (s *Convert[From, To]) Err() error {
	return s.err
}

func nextArrayElem[T any](elements []T, filter typ.Predicate[T], indexHolder *int) (T, bool) {
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
