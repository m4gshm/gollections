package set

import (
	"fmt"

	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/mutable"
	miter "github.com/m4gshm/container/mutable/iter"
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
	return WrapOrderedSet(order, uniques)
}

func NewOrderedSet[T comparable](capacity int) *OrderedSet[T] {
	return WrapOrderedSet(make([]*T, 0, capacity), make(map[T]int, capacity))
}

func WrapOrderedSet[T comparable](elements []*T, uniques map[T]int) *OrderedSet[T] {
	return &OrderedSet[T]{elements: elements, uniques: uniques}
}

type OrderedSet[T comparable] struct {
	elements   []*T
	uniques    map[T]int
	changeMark int32
}

var _ mutable.Set[any, mutable.Iterator[any]] = (*OrderedSet[any])(nil)
var _ fmt.Stringer = (*OrderedSet[any])(nil)

func (s *OrderedSet[T]) Begin() mutable.Iterator[T] {
	return s.Iter()
}

func (s *OrderedSet[T]) Iter() *RefIter[T] {
	return &RefIter[T]{s.delIter()}
}

func (s *OrderedSet[T]) delIter() *miter.Deleteable[*T] {
	return miter.NewDeleteable(&s.elements, &s.changeMark, func(ref *T) bool { return s.Delete(*ref) })
}

func (s *OrderedSet[T]) Elements() []T {
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

func (s *OrderedSet[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, typ.Iterator[T]] {
	return it.NewPipe[T](&it.RefIter[T]{Iterator: it.Filter(s.delIter(), func(ref *T) bool { return filter(*ref) })})
}

func (s *OrderedSet[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, typ.Iterator[T]] {
	return it.NewPipe[T](&it.RefIter[T]{Iterator: it.Map(s.delIter(), func(ref *T) *T {
		conv := by(*ref)
		return &conv
	})})
}

func (s *OrderedSet[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(&RefIter[T]{s.delIter()}, by)
}

func (s *OrderedSet[T]) String() string {
	return slice.ToString(s.elements)
}
