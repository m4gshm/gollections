package mutable

import (
	"errors"

	"github.com/m4gshm/gollections/typ"
)

var BadRW = errors.New("concurrent read and write")

//Vector - the container stores ordered elements, provides index access.
type Vector[t any, IT Iterator[t]] interface {
	typ.Vector[t, IT]
	Addable[t]
	Settable[int, t]
	Delete(index int) (bool, error)
}

//Set - the container provides uniqueness (does't insert duplicated values).
type Set[t comparable, IT Iterator[t]] interface {
	typ.Set[t, IT]
	Addable[t]
	Delete(...t) (bool, error)
}

//Map - the container provides access to elements by key.
type Map[k comparable, v any, IT typ.KVIterator[k, v]] interface {
	typ.Map[k, v, IT]
	Settable[k, v]
}

type Addable[T any] interface {
	Add(...T) (bool, error)
}

type Settable[k any, v any] interface {
	Set(key k, value v) (bool, error)
}

type Iterable[T any, IT typ.Iterator[T]] interface {
	Begin() IT
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
