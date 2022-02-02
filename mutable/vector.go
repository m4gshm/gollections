package mutable

import (
	"errors"
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

var BadIndex = errors.New("bad index")

func NewVector[T any](capacity int) *Vector[T] {
	return WrapVector(make([]T, 0, capacity))
}

func ToVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Copy(elements))
}

func WrapVector[T any](elements []T) *Vector[T] {
	r := &elements
	return &Vector[T]{elements: r}
}

//Vector stores ordered elements, provides index access.
type Vector[T any] struct {
	elements   *[]T
	changeMark int32
	err        error
}

var (
	_ Addable[any]       = (*Vector[any])(nil)
	_ Deleteable[int]    = (*Vector[any])(nil)
	_ Settable[int, any] = (*Vector[any])(nil)
	_ c.Vector[any]    = (*Vector[any])(nil)
	_ fmt.Stringer       = (*Vector[any])(nil)
)

func (s *Vector[T]) Begin() c.Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) BeginEdit() Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) Iter() *Iter[T] {
	return NewIter(&s.elements, &s.changeMark, s.DeleteOne)
}

func (s *Vector[T]) Collect() []T {
	return slice.Copy(*s.elements)
}

func (s *Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(*s.elements, tracker)
}

func (s *Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(*s.elements, tracker)
}

func (s *Vector[T]) For(walker func(T) error) error {
	return slice.For(*s.elements, walker)
}

func (s *Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(*s.elements, walker)
}

func (s *Vector[T]) Len() int {
	return it.GetLen(s.elements)
}

func (s *Vector[T]) Get(index int) (T, bool) {
	return slice.Get(*s.elements, index)
}

func (s *Vector[T]) Add(v ...T) (bool, error) {
	return s.AddAll(v)
}

func (s *Vector[T]) AddAll(v []T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	*s.elements = append(*s.elements, v...)
	return Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Vector[T]) AddOne(v T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	*s.elements = append(*s.elements, v)
	return Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Vector[T]) DeleteOne(index int) (bool, error) {
	_, ok, err := s.Remove(index)
	return ok, err
}

func (s *Vector[T]) Remove(index int) (T, bool, error) {
	if err := s.err; err != nil {
		var no T
		return no, false, err
	}
	if e := *s.elements; index >= 0 && index < len(e) {
		de := e[index]
		markOnStart := s.changeMark
		*s.elements = slice.Delete(index, e)
		ok, err := Commit(markOnStart, &s.changeMark, &s.err)
		return de, ok, err
	}
	var no T
	return no, false, nil
}

func (s *Vector[T]) Delete(indexes ...int) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}

	l := len(indexes)
	if l == 0 {
		return false, nil
	} else if l == 1 {
		return s.DeleteOne(indexes[0])
	}

	markOnStart := s.changeMark
	e := *s.elements
	el := len(e)

	sort.Ints(indexes)

	shift := 0
	for i := 0; i < l; i++ {
		index := indexes[i] - shift
		delAmount := 1
		if index >= 0 && index < el {
			curIndex := index
			for i < l-1 {
				nextIndex := indexes[i+1]
				if nextIndex-curIndex == 1 {
					delAmount++
					i++
					curIndex = nextIndex
				} else {
					break
				}
			}

			e = append(e[0:index], e[index+delAmount:]...)
			shift += delAmount
		}
	}
	if shift > 0 {
		*s.elements = e
		return Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *Vector[T]) Set(index int, value T) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	e := *s.elements
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
		ne := make([]T, l, c)
		copy(ne, e)
		e = ne
		*s.elements = e
	}
	e[index] = value
	return true, nil
}

func (s *Vector[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T, c.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Vector[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T, c.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *Vector[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *Vector[T]) String() string {
	return slice.ToString(*s.elements)
}
