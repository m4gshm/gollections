package iter

import (
	"errors"

	"github.com/m4gshm/container/typ"
)

const noStarted = -1

var Exhausted = errors.New("interator exhausted")

func New[T any](elements []T) *Iter[T] {
	return &Iter[T]{elements: elements, current: noStarted}
}

func NewReseteable[T any](elements []T) *Reseteable[T] {
	return &Reseteable[T]{New(elements)}
}

type Iter[T any] struct {
	elements []T
	err      error
	current  int
}

var _ typ.Iterator[interface{}] = (*Iter[interface{}])(nil)

func (s *Iter[T]) HasNext() bool {
	return hasNext(s.elements, &s.current, &s.err)
}

func (s *Iter[T]) Get() T {
	return get(s.current, s.elements, s.err)
}

func (s *Iter[T]) Err() error {
	return s.err
}

func (s *Iter[T]) Position() int {
	return s.current
}

func (s *Iter[T]) SetPosition(pos int) {
	s.current = pos
}

type Reseteable[T any] struct {
	*Iter[T]
}

var _ typ.Resetable = (*Reseteable[interface{}])(nil)

func (s *Reseteable[T]) Reset() {
	s.SetPosition(noStarted)
	s.err = nil
}

func NewDeleteable[T any](elements *[]T, changeMark *int32, del func(v T) bool) *Deleteable[T] {
	return &Deleteable[T]{elements: elements, current: noStarted, changeMark: changeMark, del: del}
}

type Deleteable[T any] struct {
	elements   *[]T
	err        error
	current    int
	changeMark *int32
	del        func(v T) bool
}

var _ typ.Iterator[any] = (*Deleteable[any])(nil)

func (i *Deleteable[T]) HasNext() bool {
	return hasNext(*i.elements, &i.current, &i.err)
}
func (i *Deleteable[T]) Get() T {
	return get(i.current, *i.elements, i.err)
}

func (i *Deleteable[T]) Delete() bool {
	pos := i.current
	if deleted := i.del(i.Get()); deleted {
		i.current = pos - 1
		return true
	}
	return false
}

func (s *Deleteable[T]) Err() error {
	return s.err
}

func hasNext[T any](elements []T, current *int, err *error) bool {
	l := len(elements)
	if l == 0 {
		*err = Exhausted
		return false
	}
	c := *current
	if c == noStarted || c < (l-1) {
		*current++
		return true
	}
	*err = Exhausted
	return false
}

func get[T any](current int, elements []T, err error) T {
	if err != nil {
		panic(err)
	}
	return (elements)[current]
}
