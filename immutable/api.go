package immutable

import (
	"github.com/m4gshm/container/typ"
)

type Vector[T any, IT typ.Iterator[T]] interface {
	typ.Vector[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
}

type Set[T any, IT typ.Iterator[T]] interface {
	typ.Set[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
}

type Map[k comparable, v any] interface {
	typ.Map[k, v]
	typ.Iterable[*typ.KV[k, v], typ.Iterator[*typ.KV[k, v]]]
}
