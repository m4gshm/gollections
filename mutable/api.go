package mutable

import (
	"errors"

	"github.com/m4gshm/container/typ"
)

var BadRW = errors.New("concurrent read and write")

//Vector - the container stores ordered elements, provides index access
type Vector[T any, IT Iterator[T]] interface {
	typ.Vector[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
	Addable[T]
	Settable[int, T]
	Delete(index int) (bool, error)
}

//Set - the container provides uniqueness (does't insert duplicated values)
type Set[T comparable, IT Iterator[T]] interface {
	typ.Set[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
	Addable[T]
	Delete(...T) (bool, error)
}

//Map - the container provides access to elements by key
type Map[k comparable, v any, IT typ.Iterator[*typ.KV[k, v]]] interface {
	typ.Map[k, v]
	typ.Iterable[*typ.KV[k, v], IT]
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
