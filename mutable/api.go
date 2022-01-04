package mutable

import (
	"github.com/m4gshm/container/typ"
)

type Vector[T any, IT typ.Iterator[T]] interface {
	typ.Vector[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
	Appendable[T]
}

type Set[T any, IT Iterator[T]] interface {
	typ.Set[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
	Appendable[T]
	Deletable[T]
}

type Map[k comparable, v any, IT typ.Iterator[*typ.KV[k, v]]] interface {
	typ.Map[k, v]
	typ.Iterable[*typ.KV[k, v], IT]
	Put(key k, value v) bool
}

type Appendable[T any] interface {
	Add(T) bool
}

type Deletable[T any] interface {
	Delete(T) bool
}

type Iterable[T any, IT typ.Iterator[T]] interface {
	Begin() IT
}

type Iterator[T any] interface {
	typ.Iterator[T]
	Delete() bool
}
