package typ

import (
	"constraints"

	"github.com/m4gshm/gollections/op"
)

//Container - base interface for container interfaces
type Container[T any, C any, IT Iterator[T]] interface {
	Walk[T]
	Collectable[C]
	Iterable[T, IT]
}

//Vector - the container stores ordered elements, provides index access
type Vector[T any, IT Iterator[T]] interface {
	Container[T, []T, IT]
	Track[T, int]
	RandomAccess[int, T]
	Transformable[T, []T, Iterator[T]]
}

//Set - the container provides uniqueness (does't insert duplicated values)
type Set[T any, IT Iterator[T]] interface {
	Container[T, []T, IT]
	Transformable[T, []T, Iterator[T]]
	Checkable[T]
}

//Map - the container provides access to elements by key
type Map[k comparable, v any, IT Iterator[*KV[k, v]]] interface {
	Container[*KV[k, v], map[k]v, Iterator[*KV[k, v]]]
	Track[v, k]
	Checkable[k]
	KeyAccess[k, v]
	MapTransformable[k, v, map[k]v]
	Keys() Container[k, []k, Iterator[k]]
	Values() Container[v, []v, Iterator[v]]
}

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element or error
	HasNext() bool
	//retrieves next element or error
	Next() (T, error)
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
	ForEach(func(element T)) error
}

//Track traverses container elements with position tracking (index, key, coordinates, etc.)
type Track[T any, P any] interface {
	TrackEach(func(position P, element T)) error
}

//Checkable container with ability to check if an element is present
type Checkable[T any] interface {
	Contains(T) bool
}

//Collectable not endless container that can be transformed to array or map of elements
type Collectable[T any] interface {
	Collect() T
}

type Transformable[T any, C any, IT Iterator[T]] interface {
	Filter(Predicate[T]) Pipe[T, C, IT]
	Map(Converter[T, T]) Pipe[T, C, IT]
	Reduce(op.Binary[T]) T
}

type Pipe[T any, C any, IT Iterator[T]] interface {
	Transformable[T, C, IT]
	Container[T, C, IT]
}

type Pipeable[T any, C any, IT Iterator[T], P Pipe[T, C, IT]] interface {
	Pipe() P
}

type MapPipe[k comparable, v any, m any] interface {
	MapTransformable[k, v, m]
	Container[*KV[k, v], m, Iterator[*KV[k, v]]]
}

type MapTransformable[k comparable, v any, m any] interface {
	Filter(Predicate[*KV[k, v]]) MapPipe[k, v, m]
	Map(Converter[*KV[k, v], *KV[k, v]]) MapPipe[k, v, m]

	FilterKey(Predicate[k]) MapPipe[k, v, m]
	MapKey(Converter[k, k]) MapPipe[k, v, m]

	FilterValue(Predicate[v]) MapPipe[k, v, m]
	MapValue(Converter[v, v]) MapPipe[k, v, m]
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
