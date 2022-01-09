package it

import (
	"errors"

	"github.com/m4gshm/container/typ"
)

const NoStarted = -1

var (
	Exhausted        = errors.New("interator exhausted")
	GetBeforeHasNext = errors.New("'Get' called before 'HasNext'")
)

func New[T any](elements []T) *Iter[T] {
	return &Iter[T]{elements: elements, current: NoStarted}
}

func NewReseteable[T any](elements []T) *Reseteable[T] {
	return &Reseteable[T]{New(elements)}
}

type Iter[T any] struct {
	elements []T
	err      error
	current  int
}

var _ typ.Iterator[any] = (*Iter[any])(nil)

func (s *Iter[T]) HasNext() bool {
	return HasNext(s.elements, &s.current, &s.err)
}

func (s *Iter[T]) Next() (T, error) {
	return Next(s.current, s.elements, s.err)
}

func (s *Iter[T]) Get() T {
	return Get(s.current, s.elements, s.err)
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
	s.SetPosition(NoStarted)
	s.err = nil
}

func HasNext[T any](elements []T, current *int, err *error) bool {
	l := len(elements)
	if l == 0 {
		*err = Exhausted
		return false
	}
	c := *current
	if c == NoStarted || c < (l-1) {
		*current++
		return true
	}
	*err = Exhausted
	return false
}

func Next[T any](current int, elements []T, err error) (T, error) {
	if err != nil {
		var no T
		return no, err
	} else if current == NoStarted {
		var no T
		return no, GetBeforeHasNext
	}
	return (elements)[current], nil
}

func Get[T any](current int, elements []T, err error) T {
	v, err := Next(current, elements, err)
	if err != nil {
		panic(err)
	}
	return v
}
