package immutable

import (
	"github.com/m4gshm/gollections/typ"
)

//Vector - the container stores ordered elements, provides index access.
type Vector[T any] typ.Vector[T, typ.Iterator[T]]

//Set - the container provides uniqueness (does't insert duplicated values).
type Set[T any] typ.Set[T, typ.Iterator[T]]

//Map - the container provides access to elements by key.
type Map[k comparable, v any] typ.Map[k, v, typ.KVIterator[k, v]]
