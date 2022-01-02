package mutable

import (
	"github.com/m4gshm/container/typ"
)

type Vector[T any] interface {
	typ.Vector[T]
	typ.Iterable[T]
	typ.Transformable[T]
	typ.Appendable[T]
}

type DelIter[T any] interface {
	typ.Iterator[T]
	Delete() bool
}

type Set[T any] interface {
	typ.Set[T]
	DelIterable[T]
	typ.Transformable[T]
	typ.Appendable[T]
	typ.Deletable[T]
}

type DelIterable[T any] interface {
	Begin() DelIter[T]
}

