package mutable

import (
	"github.com/m4gshm/container/typ"
)

type Vector[T any] interface {
	typ.Vector[T]
	typ.Iterable[T, typ.Iterator[T]]
	typ.Transformable[T]
	typ.Appendable[T]
}

type DelIter[T any] interface {
	typ.Iterator[T]
	Delete() bool
}

type Set[T any] interface {
	typ.Set[T]
	typ.Iterable[T, DelIter[T]]
	typ.Transformable[T]
	typ.Appendable[T]
	typ.Deletable[T]
}

func NewVector[T any](values ...T) Vector[T] {
	return New[T](values)
}

func NewOrderedSet[T comparable](values ...T) Set[T] {
	return NewSet(values)
}
