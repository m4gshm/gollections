package typ

import (
	"constraints"

	"github.com/m4gshm/container/op"
)

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element or error
	HasNext() bool
	//retrieves next element
	Get() T
	//retrieves error
	Err() error
}

//Iterable iterator supplier
type Iterable[T any, It Iterator[T]] interface {
	Begin() It
}

type Stream[T any] interface {
	Filter(Predicate[T]) Stream[T]
	Map(Converter[T, T]) Stream[T]
	Reduce(op.Binary[T]) T
	Iterable[T, Iterator[T]]
	Walk[T]
}

//Walk walks the elements of container
type Track[T any, P any] interface {
	ForEach(Tracker[P, T])
}

//Tracker callback function for Track interface
type Tracker[P any, T any] func(position P, value T)

type Walk[T any] interface {
	ForEach(Walker[T])
}
//Walker callback function for Walk interface
type Walker[T any] func(value T)

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
