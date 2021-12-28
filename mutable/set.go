package mutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/slice"
)

func newSet[T comparable](values []T) *OrderSet[T] {
	var (
		uniques = make(map[T]int, 0)
		order   = make([]*T, 0, 0)
	)
	pos := 0
	for _, v := range values {
		if _, ok := uniques[v]; !ok {
			vv := v
			order = append(order, &vv)
			uniques[vv] = pos
			pos++
		}
	}
	return &OrderSet[T]{elements: order, uniques: uniques}
}

type OrderSet[T comparable] struct {
	elements   []*T
	uniques    map[T]int
	changeMark int32
}

var _ Set[any] = (*OrderSet[any])(nil)
var _ fmt.Stringer = (*OrderSet[any])(nil)

func (s *OrderSet[T]) Begin() Iterator[T] {
	return &OrderIter[T]{iter.NewDeleteable(&s.elements, &s.changeMark, func(ref *T) bool {return s.Delete(*ref)})}
}

func (s *OrderSet[T]) Values() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *OrderSet[T]) Len() int {
	return len(s.elements)
}

func (s *OrderSet[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *OrderSet[T]) Add(v T) bool {
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

func (s *OrderSet[T]) Delete(v T) bool {
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

func (s *OrderSet[T]) String() string {
	return slice.ToString(s.elements)
}
