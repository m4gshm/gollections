package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func WrapKeys[k comparable](elements []k) *MapKeys[k] {
	return &MapKeys[k]{elements}
}

type MapKeys[k comparable] struct {
	elements []k
}

var _ c.Collection[int, []int, c.Iterator[int]] = (*MapKeys[int])(nil)
var _ fmt.Stringer = (*MapKeys[int])(nil)

func (s *MapKeys[k]) Begin() c.Iterator[k] {
	return s.Head()
}

func (s *MapKeys[k]) Head() *it.Iter[k] {
	return it.NewHead(s.elements)
}

func (s *MapKeys[k]) Len() int {
	return len(s.elements)
}

func (s *MapKeys[k]) Collect() []k {
	elements := s.elements
	dest := make([]k, len(elements))
	copy(dest, elements)
	return dest
}

func (s *MapKeys[k]) For(walker func(k) error) error {
	return slice.For(s.elements, walker)
}

func (s *MapKeys[k]) ForEach(walker func(k)) {
	slice.ForEach(s.elements, walker)
}

func (s *MapKeys[k]) Get(index int) (k, bool) {
	return slice.Get(s.elements, index)
}

func (s *MapKeys[t]) Filter(filter c.Predicate[t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Filter(s.Head(), filter))
}

func (s *MapKeys[t]) Map(by c.Converter[t, t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Map(s.Head(), by))
}

func (s *MapKeys[t]) Reduce(by op.Binary[t]) t {
	return it.Reduce(s.Head(), by)
}

func (s *MapKeys[k]) String() string {
	return slice.ToString(s.Collect())
}
