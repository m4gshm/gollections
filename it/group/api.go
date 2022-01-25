package group

import (
	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/typ"
)

//Of - group.Of synonym for the it.Group.
func Of[T any, K comparable, IT typ.Iterator[T]](elements IT, by typ.Converter[T, K]) typ.MapPipe[K, T, map[K][]T] {
	return it.Group[T](elements, by)
}
