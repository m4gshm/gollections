package immutable

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

type Vector[T any] interface {
	typ.Vector[T]
	typ.Iterable[T, typ.Iterator[T]]
	typ.Transformable[T]
}

type Set[T any, IT typ.Iterator[T]] interface {
	typ.Set[T]
	typ.Iterable[T, IT]
	typ.Transformable[T]
}

type Map[k comparable, v any, IT typ.Iterator[*typ.KV[k, v]]] interface {
	typ.Map[k, v, IT]
}

func NewVector[T any](values ...T) Vector[T] {
	return New(values)
}

func NewOrderedSet[T comparable](values ...T) Set[T, typ.Iterator[T]] {
	return NewSet(values)
}

func NewOrderedMap[k comparable, v any](values ...*typ.KV[k, v]) typ.Map[k, v, *iter.OrderedKVIter[k, v]] {
	return newMap(values)
}
