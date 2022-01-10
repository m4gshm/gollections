package immutable

import (
	"github.com/m4gshm/gollections/typ"
)

//Vector - the container stores ordered elements, provides index access
type Vector[T any, IT typ.Iterator[T]] interface {
	typ.Vector[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
}

//Set - the container provides uniqueness (does't insert duplicated values)
type Set[T any, IT typ.Iterator[T]] interface {
	typ.Set[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
}

//Map - the container provides access to elements by key
type Map[k comparable, v any] interface {
	typ.Map[k, v]
	typ.Iterable[*typ.KV[k, v], typ.Iterator[*typ.KV[k, v]]]
}
