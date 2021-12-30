package mutable

import "github.com/m4gshm/container/typ"


type Iterator[T any] interface {
	typ.Iterator[T]
	Delete() bool
}

type Set[T any, It Iterator[T]] interface {
	typ.Walk[T]
	typ.Iterable[T, It]
	typ.Container[[]T]
	typ.Measureable[int]
	typ.Checkable[T]
	typ.Appendable[T]
	typ.Deletable[T]
}


func NewOrderedSet[T comparable](values ...T) Set[T, *OrderIter[T]] {
	return NewSet(values)
}