package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
)

//NewPipe returns the Pipe based on iterable elements.
func NewPipe[t any, it c.Iterator[t]](iter it) *IterPipe[t] {
	return &IterPipe[t]{it: iter}
}

type IterPipe[t any] struct {
	it       c.Iterator[t]
	elements []t
}

var _ c.Pipe[any, []any, c.Iterator[any]] = (*IterPipe[any])(nil)

func (s *IterPipe[t]) Filter(fit c.Predicate[t]) c.Pipe[t, []t, c.Iterator[t]] {
	return NewPipe[t](Filter(s.it, fit))
}

func (s *IterPipe[t]) Map(by c.Converter[t, t]) c.Pipe[t, []t, c.Iterator[t]] {
	return NewPipe[t](Map(s.it, by))
}

func (s *IterPipe[t]) ForEach(walker func(t)) {
	ForEach(s.it, walker)
}

func (s *IterPipe[t]) For(walker func(t) error) error {
	return For(s.it, walker)
}

func (s *IterPipe[t]) Reduce(by op.Binary[t]) t {
	return Reduce(s.it, by)
}

func (s *IterPipe[t]) Begin() c.Iterator[t] {
	return s.it
}

func (s *IterPipe[t]) Collect() []t {
	e := s.elements
	if e == nil {
		e = make([]t, 0)
		it := s.it
		for it.HasNext() {
			e = append(e, it.Next())
		}
		s.elements = e
	}
	return e
}
