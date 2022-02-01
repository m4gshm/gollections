//Package mutable provides implementations of mutable containers.
package mutable

import (
	"errors"

	"github.com/m4gshm/gollections/typ"
)

var BadRW = errors.New("concurrent read and write")

//Map provides access to elements by key.
type Map[k comparable, v any] interface {
	typ.Map[k, v]
	Settable[k, v]
}

type Addable[T any] interface {
	Add(...T) (bool, error)
}

type Settable[k any, v any] interface {
	Set(key k, value v) (bool, error)
}

type Deleteable[k any] interface {
	Delete(...k) (bool, error)
}

type Removable[k any, v any] interface {
	Remove(k) (v, bool, error)
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
