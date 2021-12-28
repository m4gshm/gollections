package typ

import (
	"constraints"
)

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element
	HasNext() bool
	//retrieves next element
	Get() T
}

type Iterable[T any, It Iterator[T]] interface {
	Begin() It
}

type Walker[P any, T any] func(position P, value T)
type Walk[T any, P any] interface {
	ForEach(Walker[P, T])
}

type Container[T any] interface {
	Values() T
}

type Measureable[L constraints.Integer] interface {
	Len() L
}

type Checkable[T any] interface {
	Contains(T) bool
}

type Resetable interface {
	Reset()
}

type Appendable[T any] interface {
	Add(T) bool
}

type Deletable[T any] interface {
	Delete(T) bool
}

type Access[K any, V any] interface {
	Get(K) V
}

type RandomAccess[K constraints.Integer, V any] interface {
	Access[K, V]
}

type KeyAccess[K comparable, V any] interface {
	Access[K, V]
}
