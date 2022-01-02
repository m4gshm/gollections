package set

import (
	"fmt"

	"github.com/m4gshm/container/immutable/set"
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func ToOrderedSet[T comparable](elements []T) *OrderedSet[T] {
	var (
		uniques = make(map[T]int, 0)
		order   = make([]*T, 0, 0)
	)
	pos := 0
	for _, v := range elements {
		if _, ok := uniques[v]; !ok {
			vv := v
			order = append(order, &vv)
			uniques[vv] = pos
			pos++
		}
	}
	return &OrderedSet[T]{elements: order, uniques: uniques}
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{elements: make([]*T, 0), uniques: map[T]int{}}
}

type OrderedSet[T comparable] struct {
	elements   []*T
	uniques    map[T]int
	changeMark int32
}

var _ mutable.Set[any] = (*OrderedSet[any])(nil)
var _ fmt.Stringer = (*OrderedSet[any])(nil)

func (s *OrderedSet[T]) Begin() mutable.DelIter[T] {
	return s.Iter()
}

func (s *OrderedSet[T]) Iter() *OrderIter[T] {
	return &OrderIter[T]{s.delIter()}
}

func (s *OrderedSet[T]) delIter() *iter.Deleteable[*T] {
	return iter.NewDeleteable(&s.elements, &s.changeMark, func(ref *T) bool { return s.Delete(*ref) })
}

func (s *OrderedSet[T]) Values() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *OrderedSet[T]) ForEach(walker func(T)) {
	for _, e := range s.elements {
		walker(*e)
	}
}

func (s *OrderedSet[T]) Len() int {
	return len(s.elements)
}

func (s *OrderedSet[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *OrderedSet[T]) Add(v T) bool {
	markOnStart := s.changeMark
	u := s.uniques
	if _, ok := u[v]; !ok {
		e := s.elements
		u[v] = len(e)
		s.elements = append(e, &v)
		markOnFinish := s.changeMark
		if markOnFinish != markOnStart {
			panic("concurrent ordered map read and write")
		}
		s.changeMark++
		return true
	}
	return false
}

func (s *OrderedSet[T]) Delete(v T) bool {
	markOnStart := s.changeMark
	u := s.uniques
	changeMark := s.changeMark
	if pos, ok := u[v]; ok {
		delete(u, v)
		e := s.elements
		ne := append(e[0:pos], e[pos+1:]...)
		for i := pos; i < len(ne); i++ {
			u[*ne[i]]--
		}
		s.elements = ne
		markOnFinish := s.changeMark
		if markOnFinish != markOnStart {
			panic("concurrent ordered map read and write")
		}
		changeMark++
		return true
	}
	return false
}

func (s *OrderedSet[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T] {
	return iter.NewPipe[T](&set.OrderIter[T]{Iterator: iter.Filter(s.delIter(), func(ref *T) bool { return filter(*ref) })})
}

func (s *OrderedSet[T]) Map(by typ.Converter[T, T]) typ.Pipe[T] {
	return iter.NewPipe[T](&set.OrderIter[T]{Iterator: iter.Map(s.delIter(), func(ref *T) *T {
		conv := by(*ref)
		return &conv
	})})
}

func (s *OrderedSet[T]) Reduce(by op.Binary[T]) T {
	return iter.Reduce(&OrderIter[T]{s.delIter()}, by)
}

func (s *OrderedSet[T]) String() string {
	return slice.ToString(s.elements)
}
