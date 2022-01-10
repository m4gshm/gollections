package vector

import (
	"errors"
	"fmt"

	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/slice"
)

var BadIndex = errors.New("bad index")

func Create[T any](capacity int) *Vector[T] {
	return Wrap(make([]T, 0, capacity))
}

func Convert[T any](elements []T) *Vector[T] {
	return Wrap(slice.Copy(elements))
}

func Wrap[T any](elements []T) *Vector[T] {
	r := &elements
	return &Vector[T]{Vector: vector.Wrap(r), elements: &r}
}

type Vector[t any /*if replaces generic type by 'T' it raises compile-time error 'type parameter bound more than once'*/] struct {
	*vector.Vector[t]
	elements   **[]t
	changeMark int32
	err        error
}

var _ mutable.Vector[any, mutable.Iterator[any]] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[t]) Begin() mutable.Iterator[t] {
	return s.Iter()
}

func (s *Vector[t]) Iter() *Iter[t] {
	return NewIter(s.elements, &s.changeMark, s.Delete)
}

func (s *Vector[t]) Add(v ...t) (bool, error) {
	return s.AddAll(v)
}

func (s *Vector[t]) AddAll(v []t) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	**s.elements = append(**s.elements, v...)
	return mutable.Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Vector[t]) AddOne(v t) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	**s.elements = append(**s.elements, v)
	return mutable.Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Vector[t]) Delete(index int) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	e := **s.elements
	if index >= 0 && index < len(e) {
		markOnStart := s.changeMark
		**s.elements = slice.Delete(index, e)
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *Vector[t]) Set(index int, value t) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	e := **s.elements
	if index < 0 {
		return false, fmt.Errorf("%w: %d", BadIndex, index)
	}
	l := len(e)
	if index >= l {
		c := l * 2
		l := index + 1
		if l > c {
			c = l
		}
		ne := make([]t, l, c)
		copy(ne, e)
		e = ne
		**s.elements = e
	}
	e[index] = value
	return true, nil
}
