package immutable

import (
	"github.com/m4gshm/container/typ"
)

type Vec[T any] interface {
	typ.Vector[T]
	typ.Iterable[T]
	typ.Transformable[T]
}

type Set[T any] interface {
	typ.Set[T]
	typ.Iterable[T]
	typ.Transformable[T]
}

type Map[k comparable, v any] interface {
	typ.Map[k, v]
}

