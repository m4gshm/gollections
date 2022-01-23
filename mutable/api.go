//Package mutable provides implementations of mutable containers.
package mutable

import (
	"errors"

	"github.com/m4gshm/gollections/typ"
)

var BadRW = errors.New("concurrent read and write")

//Vector stores ordered elements, provides index access.
type Vector[t any] interface {
	typ.Vector[t, typ.Iterator[t]]
	Addable[t]
	Settable[int, t]
	BeginEdit() Iterator[t]
	Delete(index int) (bool, error)
}

//Set provides uniqueness (does't insert duplicated values).
type Set[t comparable] interface {
	typ.Set[t, typ.Iterator[t]]
	Addable[t]
	BeginEdit() Iterator[t]
	Delete(...t) (bool, error)
}

//Map provides access to elements by key.
type Map[k comparable, v any] interface {
	typ.Map[k, v, typ.KVIterator[k, v]]
	Settable[k, v]
}

type Addable[T any] interface {
	Add(...T) (bool, error)
}

type Settable[k any, v any] interface {
	Set(key k, value v) (bool, error)
}

type Iterator[T any] interface {
	typ.Iterator[T]
	Delete() (bool, error)
}

func Commit(markOnStart int32, mark *int32, err *error) (bool, error) {
	markOnFinish := *mark
	if markOnFinish != markOnStart {
		e := BadRW
		*err = e
		return false, e
	}
	(*mark)++
	return true, nil
}
