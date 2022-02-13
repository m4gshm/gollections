//Package c provides common types of containers, utility types and functions.
package c

import (
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
)

// Container is the base interface of implementations.
// Where:
// 		Collection - a collection type (may be slice, map or chain)
// 		IT - a iterator type  (Iterator, KVIterator)
type Container[Collection any, IT Iter] interface {
	Collect() Collection
	Iterable[IT]
}

//Collection is the interface of a finite-size container.
// Where:
//		T - any arbitrary type
// 		Collection - a collection type (may be slice, map or chain)
// 		IT - a iterator type  (Iterator, KVIterator)
type Collection[T any, Collection any, IT Iter] interface {
	Container[Collection, IT]
	Walk[T]
	WalkEach[T]
	Len() int
	IsEmpty() bool
}

//Vector is the interface of a container stores ordered elements, provides index access.
type Vector[T any] interface {
	Collection[T, []T, Iterator[T]]
	Track[T, int]
	TrackEach[T, int]
	Access[int, T]
	Transformable[T, []T]
}

//Set is the interface of a container provides uniqueness (does't insert duplicated values).
type Set[T any] interface {
	Collection[T, []T, Iterator[T]]
	Transformable[T, []T]
	Checkable[T]
}

//Map is the interface of a container provides access to elements by key.
type Map[K comparable, V any] interface {
	Collection[*map_.KV[K, V], map[K]V, KVIterator[K, V]]
	Track[V, K]
	TrackEach[V, K]
	Checkable[K]
	Access[K, V]
	MapTransformable[K, V, map[K]V]
	Keys() Collection[K, []K, Iterator[K]]
	Values() Collection[V, []V, Iterator[V]]
}

//Iter is the base for Iterator and KVIterator.
type Iter interface {
	//checks ability on next element or error.
	HasNext() bool
}

//Iterator is the interface provides iterate over elements of a collection.
type Iterator[T any] interface {
	Iter
	//retrieves next element or zero value if no more elements
	//must be called only after HasNext
	Get() T
}

//KVIterator is the interface provides iterate over all key/value pair of a map.
type KVIterator[K, V any] interface {
	Iter
	//retrieves next elements or zero values if no more elements
	//must be called only after HasNext
	Get() (K, V)
}

//Iterable is an iterator supplier interface
type Iterable[IT Iter] interface {
	Begin() IT
}

//Walk touches elements of the collection.
type Walk[IT any] interface {
	For(func(element IT) error) error
}

//WalkEach touches all elements of the collection without error checking
type WalkEach[T any] interface {
	ForEach(func(element T))
}

//Track traverses container elements with position tracking (index, key, coordinates, etc.)
type Track[T any, P any] interface {
	Track(func(position P, element T) error) error
}

//TrackEach traverses container elements with position tracking (index, key, coordinates, etc.) without error checking
type TrackEach[T any, P any] interface {
	TrackEach(func(position P, element T))
}

//Checkable container with ability to check if an element is present.
type Checkable[T any] interface {
	Contains(T) bool
}

//Transformable is the interface that provides limited kit of container transformation methods.
//The full kit of transformer functions are in the package 'c'
type Transformable[T any, Collection any] interface {
	Filter(Predicate[T]) Pipe[T, Collection]
	Map(Converter[T, T]) Pipe[T, Collection]
}

//Pipe extends Transformable by finalize methods like ForEach, Collect or Reduce.
type Pipe[T any, Collection any] interface {
	Transformable[T, Collection]
	Container[Collection, Iterator[T]]
	Walk[T]
	WalkEach[T]
	Reduce(op.Binary[T]) T
}

//MapTransformable is the interface that provides limited kit of map transformation methods.
//The full kit of transformer functions are in the package 'c/map_'
type MapTransformable[K comparable, V any, Map any] interface {
	Filter(BiPredicate[K, V]) MapPipe[K, V, Map]
	Map(BiConverter[K, V, K, V]) MapPipe[K, V, Map]

	FilterKey(Predicate[K]) MapPipe[K, V, Map]
	MapKey(Converter[K, K]) MapPipe[K, V, Map]

	FilterValue(Predicate[V]) MapPipe[K, V, Map]
	MapValue(Converter[V, V]) MapPipe[K, V, Map]
}

//MapPipe extends MapTransformable by finalize methods like ForEach, Collect or Reduce.
type MapPipe[K comparable, V any, Map any] interface {
	MapTransformable[K, V, Map]
	Container[Map, KVIterator[K, V]]
}

//Access is the interface that provides access to an element by its pointer (index, key, coordinate, etc.)
// Where:
//		P - a type of pointer to a value (index, map key, coordinates)
// 		V - any arbitrary type of the value
type Access[P any, V any] interface {
	Get(P) (V, bool)
}
