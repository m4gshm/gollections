package oset

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
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
	err        error
}

var _ mutable.Set[any] = (*OrderedSet[any])(nil)
var _ typ.Set[any, typ.Iterator[any]] = (*OrderedSet[any])(nil)
var _ fmt.Stringer = (*OrderedSet[any])(nil)

func (s *OrderedSet[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *OrderedSet[T]) BeginEdit() mutable.Iterator[T] {
	return s.Iter()
}

func (s *OrderedSet[T]) Iter() *Iter[T] {
	return NewIter(&s.elements, &s.changeMark, s.DeleteOne)
}

func (s *OrderedSet[T]) Collect() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *OrderedSet[T]) For(walker func(T) error) error {
	return slice.ForRefs(s.elements, walker)
}

func (s *OrderedSet[T]) ForEach(walker func(T)) {
	slice.ForEachRef(s.elements, walker)
}

func (s *OrderedSet[T]) Len() int {
	return len(s.elements)
}

func (s *OrderedSet[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *OrderedSet[T]) Add(elements ...T) (bool, error) {
	return s.AddAll(elements)
}

func (s *OrderedSet[T]) AddAll(elements []T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	u := s.uniques
	result := false
	for i := range elements {
		markOnStart := s.changeMark
		v := elements[i]
		if _, ok := u[v]; !ok {
			e := s.elements
			u[v] = len(e)
			s.elements = append(e, &v)
			cmt, err := mutable.Commit(markOnStart, &s.changeMark, &s.err)
			if err != nil {
				return false, err
			}
			result = result || cmt
		}
	}
	return result, nil
}

func (s *OrderedSet[T]) AddOne(v T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	u := s.uniques
	if _, ok := u[v]; !ok {
		e := s.elements
		u[v] = len(e)
		s.elements = append(e, &v)
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *OrderedSet[T]) Delete(elements ...T) (bool, error) {
	return s.DeleteAll(elements)
}

func (s *OrderedSet[T]) DeleteAll(elements []T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	u := s.uniques
	result := false
	for i := range elements {
		v := elements[i]
		if pos, ok := u[v]; ok {
			markOnStart := s.changeMark
			delete(u, v)
			//todo: need optimize
			e := s.elements
			ne := slice.Delete(pos, e)
			for i := pos; i < len(ne); i++ {
				u[*ne[i]]--
			}
			s.elements = ne
			ok, err := mutable.Commit(markOnStart, &s.changeMark, &s.err)
			if err != nil {
				return false, err
			}
			result = result || ok
		}
	}
	return result, nil
}

func (s *OrderedSet[T]) DeleteOne(v T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	u := s.uniques
	if pos, ok := u[v]; ok {
		markOnStart := s.changeMark
		delete(u, v)
		//todo: need optimize
		e := s.elements
		ne := slice.Delete(pos, e)
		for i := pos; i < len(ne); i++ {
			u[*ne[i]]--
		}
		s.elements = ne
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *OrderedSet[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *OrderedSet[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *OrderedSet[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *OrderedSet[T]) String() string {
	return slice.ToString(s.elements)
}
