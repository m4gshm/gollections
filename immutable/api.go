//Package immutable provides implementations of constant containers.
package immutable

import (
	"github.com/m4gshm/gollections/typ"
)

//Vector stores ordered elements, provides index access.
type Vector[T any] typ.Vector[T, typ.Iterator[T]]

//Set provides uniqueness (does't insert duplicated values).
type Set[T any] typ.Set[T, typ.Iterator[T]]

//Map provides access to elements by key.
type Map[k comparable, v any] typ.Map[k, v, typ.KVIterator[k, v]]
