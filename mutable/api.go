//Package mutable provides implementations of mutable containers.
package mutable

import (
	"github.com/m4gshm/gollections/c"
)

type Addable[T any] interface {
	Add(...T) bool
}

type Settable[k any, v any] interface {
	Set(key k, value v) bool
}

type Deleteable[k any] interface {
	Delete(...k) bool
}

type Removable[k any, v any] interface {
	Remove(k) (v, bool)
}

type Iterator[T any] interface {
	c.Iterator[T]
	Delete() bool
}

