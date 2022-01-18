package typ

import (
	"constraints"

	"github.com/m4gshm/gollections/op"
)

//Container - base interface for container interfaces
type Container[C any, IT Iter] interface {
	Collectable[C]
	Iterable[IT]
}

//Vector - the container stores ordered elements, provides index access
type Vector[T any, IT Iterator[T]] interface {
	Walk[T]
	Container[[]T, IT]
	Track[T, int]
	RandomAccess[int, T]
	Transformable[T, []T, Iterator[T]]
	Len() int
}

//Set - the container provides uniqueness (does't insert duplicated values)
type Set[T any, IT Iterator[T]] interface {
	Iterable[IT]
	Walk[T]
	Container[[]T, IT]
	Transformable[T, []T, Iterator[T]]
	Checkable[T]
	Len() int
}

//Map - the container provides access to elements by key
type Map[k comparable, v any, IT KVIterator[k, v]] interface {
	Container[map[k]v, IT]
	Iterable[IT]
	Track[v, k]
	Checkable[k]
	KeyAccess[k, v]
	MapTransformable[k, v, map[k]v]
	Keys() Container[[]k, Iterator[k]]
	Values() Container[[]v, Iterator[v]]
	Len() int
}

type Iter interface {
	//checks ability on next element or error
	HasNext() bool
}

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	Iter
	//retrieves next element or error
	//must be called only after HasNext
	Get() (T, error)
}

type KVIterator[k, v any] interface {
	Iter
	//retrieves next element or error
	//must be called only after HasNext
	Get() (k, v, error)
}

//Resetable an object with resettable state (e.g. slice based iterator)
type Resetable interface {
	Reset()
}

//Iterable iterator supplier
type Iterable[IT Iter] interface {
	Begin() IT
}

//Walk touches all elements of the collection
type Walk[T any] interface {
	ForEach(func(element T)) error
	For(func(element T) error) error
}

//Track traverses container elements with position tracking (index, key, coordinates, etc.)
type Track[T any, P any] interface {
	TrackEach(func(position P, element T)) error
	Track(func(position P, element T) error) error
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
	Container[C, IT]
	Walk[T]
}

type Pipeable[T any, C any, IT Iterator[T], P Pipe[T, C, IT]] interface {
	Pipe() P
}

type MapPipe[k comparable, v any, m any] interface {
	MapTransformable[k, v, m]
	Container[m, KVIterator[k, v]]
}

type MapTransformable[k comparable, v any, m any] interface {
	Filter(BiPredicate[k, v]) MapPipe[k, v, m]
	Map(BiConverter[k, v, k, v]) MapPipe[k, v, m]

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
