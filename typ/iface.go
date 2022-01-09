package typ

import (
	"constraints"

	"github.com/m4gshm/container/op"
)

//Container - base interface for container interfaces
type Container[T any, L constraints.Integer, IT Iterator[T]] interface {
	Walk[T]
	Finite[[]T, L]
	Iterable[T, IT]
}

//Vector - the container stores ordered elements, provides index access
type Vector[T any, IT Iterator[T]] interface {
	Container[T, int, IT]
	RandomAccess[int, T]
}

//Set - the container provides uniqueness (does't insert duplicated values)
type Set[T any, IT Iterator[T]] interface {
	Container[T, int, IT]
	Checkable[T]
}

//Map - the container provides access to elements by key
type Map[k comparable, v any] interface {
	Track[v, k]
	Finite[map[k]v, int]
	Checkable[k]
	KeyAccess[k, v]
	Keys() Container[k, int, Iterator[k]]
	Values() Container[v, int, Iterator[v]]
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
type Iterable[T any, IT Iterator[T]] interface {
	Begin() IT
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
	Elements() T
	Len() L
}

type Transformable[T any, IT Iterator[T]] interface {
	Filter(Predicate[T]) Pipe[T, IT]
	Map(Converter[T, T]) Pipe[T, IT]
	Reduce(op.Binary[T]) T
}

type Pipe[T any, IT Iterator[T]] interface {
	Transformable[T, IT]
	Iterable[T, IT]
	Walk[T]
}

type Pipeable[T any, IT Iterator[T]] interface {
	Pipe() Pipe[T, IT]
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
