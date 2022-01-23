package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/typ"
)

func Of[T any, K comparable, IT typ.Iterable[typ.Iterator[T]]](elements IT, by typ.Converter[T, K]) typ.MapPipe[K, T, map[K][]T] {
	return c.Group[T, K, typ.Iterable[typ.Iterator[T]], typ.Iterator[T]](elements, by)
}
