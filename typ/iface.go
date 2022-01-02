package typ

import (
	"constraints"

	"github.com/m4gshm/container/op"
)

/*
 *  Common interfaces
 */

type Container[T any, L constraints.Integer] interface {
	Walk[T]
	Finite[[]T, L]
}

type Vector[T any] interface {
	Container[T, int]
	RandomAccess[int, T]
}

type Set[T any] interface {
	Container[T, int]
	Checkable[T]
}

type Map[k comparable, v any] interface {
	Track[v, k]
	Iterable[*KV[k, v]]
	Finite[map[k]v, int]
	Checkable[k]
	KeyAccess[k, v]
}

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element or error
	HasNext() bool
	//retrieves next element
	Get() T
	//retrieves error
	Err() error
}

//Resetable an object with resettable state (e.g. slice based iterator)
type Resetable interface {
	Reset()
}

//Iterable iterator supplier
type Iterable[T any] interface {
	Begin() Iterator[T]
}

//Walk touches all elements of the collection
type Walk[T any] interface {
	ForEach(func(element T))
}

//Track traverses container elements with position tracking (index, key, coordinates, etc.)
type Track[T any, P any] interface {
	ForEach(func(position P, element T))
}

//Checkable container with ability to check if an element is present
type Checkable[T any] interface {
	Contains(T) bool
}

//Finite not endless container that can be transformed to array or map of elements
type Finite[T any, L constraints.Integer] interface {
	Values() T
	Len() L
}

type Transformable[T any] interface {
	Filter(Predicate[T]) Pipe[T]
	Map(Converter[T, T]) Pipe[T]
	Reduce(op.Binary[T]) T
}

type Pipe[T any] interface {
	Transformable[T]
	Iterable[T]
	Walk[T]
}

type Pipeable[T any] interface {
	Pipe() Pipe[T]
}

type Appendable[T any] interface {
	Add(T) bool
}

type Deletable[T any] interface {
	Delete(T) bool
}

type Access[K any, V any] interface {
	Get(K) (V, bool)
}

type RandomAccess[K constraints.Integer, V any] interface {
	Access[K, V]
}

type KeyAccess[K comparable, V any] interface {
	Access[K, V]
}
