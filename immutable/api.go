package immutable

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

type Set[T any] interface {
	typ.Walk[T, typ.Iterator[T]]
	typ.Container[[]T]
	typ.Measureable[int]
	typ.Checkable[T]
}

type Map[k comparable, v any, It typ.Iterator[*typ.KV[k, v]]] interface {
	typ.Walk[*typ.KV[k, v], It]
	typ.Container[map[k]v]
	typ.Measureable[int]
	typ.Checkable[k]
}

func NewOrderedSet[T comparable](values ...T) Set[T] {
	return newSet[T](values)
}

func NewOrderedMap[k comparable, v any](values ...*typ.KV[k, v]) Map[k, v, *iter.OrderedKVIter[k, v]] {
	return newMap[k,v](values)
}
