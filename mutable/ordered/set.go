package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func ToSet[T comparable](elements []T) *Set[T] {
	var (
		uniques = make(map[T]int, 0)
		order   = make([]T, 0, 0)
	)
	pos := 0
	for _, v := range elements {
		if _, ok := uniques[v]; !ok {
			order = append(order, v)
			uniques[v] = pos
			pos++
		}
	}
	return WrapSet(order, uniques)
}

func NewSet[T comparable](capacity int) *Set[T] {
	return WrapSet(make([]T, 0, capacity), make(map[T]int, capacity))
}

func WrapSet[T comparable](elements []T, uniques map[T]int) *Set[T] {
	return &Set[T]{elements: elements, uniques: uniques}
}

type Set[T comparable] struct {
	elements   []T
	uniques    map[T]int
	changeMark int32
	err        error
}

var (
	_ mutable.Addable[any]    = (*Set[any])(nil)
	_ mutable.Deleteable[any] = (*Set[any])(nil)
	_ typ.Set[any]            = (*Set[any])(nil)
	_ fmt.Stringer            = (*Set[any])(nil)
)

func (s *Set[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Set[T]) BeginEdit() mutable.Iterator[T] {
	return s.Iter()
}

func (s *Set[T]) Iter() *SetIter[T] {
	return NewSetIter(&s.elements, &s.changeMark, s.DeleteOne)
}

func (s *Set[T]) Collect() []T {
	return slice.Copy(s.elements)
}

func (s *Set[T]) For(walker func(T) error) error {
	return slice.For(s.elements, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	slice.ForEach(s.elements, walker)
}

func (s *Set[T]) Len() int {
	return len(s.elements)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *Set[T]) Add(elements ...T) (bool, error) {
	return s.AddAll(elements)
}

func (s *Set[T]) AddAll(elements []T) (bool, error) {
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
			s.elements = append(e, v)
			cmt, err := mutable.Commit(markOnStart, &s.changeMark, &s.err)
			if err != nil {
				return false, err
			}
			result = result || cmt
		}
	}
	return result, nil
}

func (s *Set[T]) AddOne(v T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	u := s.uniques
	if _, ok := u[v]; !ok {
		e := s.elements
		u[v] = len(e)
		s.elements = append(e, v)
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *Set[T]) Delete(elements ...T) (bool, error) {
	return s.DeleteAll(elements)
}

func (s *Set[T]) DeleteAll(elements []T) (bool, error) {
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
				u[ne[i]]--
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

func (s *Set[T]) DeleteOne(v T) (bool, error) {
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
			u[ne[i]]--
		}
		s.elements = ne
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *Set[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Set[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *Set[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.elements)
}
