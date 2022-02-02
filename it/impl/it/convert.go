package it

import "github.com/m4gshm/gollections/c"

type ConvertFit[From, To any] struct {
	Iter    c.Iterator[From]
	By      c.Converter[From, To]
	Fit     c.Predicate[From]
	current To
	err     error
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

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

func (s *ConvertFit[From, To]) Get() (To, error) {
	return s.current, s.err
}

func (s *ConvertFit[From, To]) Next() To {
	return Next[To](s)
}

type Convert[From, To any, IT c.Iterator[From], C c.Converter[From, To]] struct {
	Iter IT
	By   C
}

var _ c.Iterator[any] = (*Convert[any, any, c.Iterator[any], c.Converter[any, any]])(nil)

func (s *Convert[From, To, IT, C]) HasNext() bool {
	return s.Iter.HasNext()
}

func (s *Convert[From, To, IT, C]) Get() (To, error) {
	v, err := s.Iter.Get()
	if err != nil {
		var no To
		return no, err
	}
	return s.By(v), nil
}

func (s *Convert[From, To, IT, C]) Next() To {
	v, err := s.Iter.Get()
	if err != nil {
		panic(err)
	}
	return s.By(v)
}

type ConvertKV[k, v any, IT c.KVIterator[k, v], k2, v2 any, C c.BiConverter[k, v, k2, v2]] struct {
	Iter IT
	By   C
}

var _ c.KVIterator[any, any] = (*ConvertKV[any, any, c.KVIterator[any, any], any, any, c.BiConverter[any, any, any, any]])(nil)

func (s *ConvertKV[k, v, IT, k1, v2, C]) HasNext() bool {
	return s.Iter.HasNext()
}

func (s *ConvertKV[k, v, IT, k2, v2, C]) Get() (k2, v2, error) {
	key, val, err := s.Iter.Get()
	if err != nil {
		var key k2
		var val v2
		return key, val, err
	}
	key1, val2 := s.By(key, val)
	return key1, val2, nil
}

func (s *ConvertKV[k, v, IT, k2, v2, C]) Next() (k2, v2) {
	return NextKV[k2, v2](s)
}
