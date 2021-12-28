package immutable

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

type Set[T any, It typ.Iterator[T]] interface {
	typ.Walk[T, int]
	typ.Iterable[T, It]
	typ.Container[[]T]
	typ.Measureable[int]
	typ.Checkable[T]
}

type Map[k comparable, v any, It typ.Iterator[*typ.KV[k, v]]] interface {
	typ.Walk[v, k]
	typ.Iterable[*typ.KV[k, v], It]
	typ.Container[map[k]v]
	typ.Measureable[int]
	typ.Checkable[k]
}

func NewOrderedSet[T comparable](values ...T) Set[T, *OrderIter[T]] {
	return NewSet(values)
}

func NewOrderedMap[k comparable, v any](values ...*typ.KV[k, v]) Map[k, v, *iter.OrderedKVIter[k, v]] {
	return newMap(values)
}
