package it

import (
	"github.com/m4gshm/gollections/typ"
)

type ConvertFit[From, To any] struct {
	Iter    typ.Iterator[From]
	By      typ.Converter[From, To]
	Fit     typ.Predicate[From]
	current To
	err     error
}

var _ typ.Iterator[interface{}] = (*ConvertFit[interface{}, interface{}])(nil)

func (s *ConvertFit[From, To]) HasNext() bool {
	if v, ok, err := nextFiltered(s.Iter, s.Fit); err != nil {
		s.err = err
		return true
	} else if ok {
		s.current = s.By(v)
		return true
	}
	s.err = Exhausted
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
	Iter typ.Iterator[From]
	By   typ.Converter[From, To]
}

var _ typ.Iterator[interface{}] = (*Convert[interface{}, interface{}])(nil)

func (s *Convert[From, To]) HasNext() bool {
	return s.Iter.HasNext()
}

func (s *Convert[From, To]) Next() (To, error) {
	v, err := s.Iter.Next()
	if err != nil {
		var no To
		return no, err
	}
	return s.By(v), nil
}
