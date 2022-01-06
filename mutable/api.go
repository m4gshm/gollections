package mutable

import (
	"errors"

	"github.com/m4gshm/container/typ"
)

var BadRW = errors.New("concurrent read and write")

type Vector[T any, IT typ.Iterator[T]] interface {
	typ.Vector[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
	Appendable[T]
}

type Set[T any, IT Iterator[T]] interface {
	typ.Set[T, IT]
	typ.Transformable[T, typ.Iterator[T]]
	Appendable[T]
	Deletable[T]
}

type Map[k comparable, v any, IT typ.Iterator[*typ.KV[k, v]]] interface {
	typ.Map[k, v]
	typ.Iterable[*typ.KV[k, v], IT]
	Put(key k, value v) bool
}

type Appendable[T any] interface {
	Add(...T) (bool, error)
}

type Deletable[T any] interface {
	Delete(...T) (bool, error)
}

type Iterable[T any, IT typ.Iterator[T]] interface {
	Begin() IT
}

type Iterator[T any] interface {
	typ.Iterator[T]
	Delete() (bool, error)
}

func Commit(markOnStart int32, mark *int32) (bool, error) {
	markOnFinish := *mark
	if markOnFinish != markOnStart {
		return false, BadRW
	}
	(*mark)++
	return true, nil
}
