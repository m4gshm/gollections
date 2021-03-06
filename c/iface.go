//Package c provides common types of containers, utility types and functions.
package c

import (
	"github.com/m4gshm/gollections/op"
)

//Container is the base interface of implementations.
//Where:
//		Collection - a collection type (may be slice, map or chain)
// 		IT - a iterator type  (Iterator, KVIterator)
type Container[Collection any, IT any] interface {
	Collect() Collection
	Iterable[IT]
}

//Collection is the interface of a finite-size container.
//Where:
//		T - any arbitrary type
// 		Collection - a collection type (may be slice, map or chain)
// 		IT - a iterator type  (Iterator, KVIterator)
type Collection[T any, Collection any, IT any] interface {
	Container[Collection, IT]
	Walk[T]
	WalkEach[T]
	Len() int
	IsEmpty() bool
}

//Vector - collection interface that provides elements order and access by index to the elements.
type Vector[T any] interface {
	Collection[T, []T, Iterator[T]]
	Track[T, int]
	TrackEach[T, int]
	Access[int, T]
	Transformable[T, []T]
}

//Set - collection interface that ensures the uniqueness of elements (does not insert duplicate values).
type Set[T any] interface {
	Collection[T, []T, Iterator[T]]
	Transformable[T, []T]
	Checkable[T]
}

//Map - collection interface that stores key/value pairs and provide access to an element by its key.
type Map[K comparable, V any] interface {
	Collection[KV[K, V], map[K]V, KVIterator[K, V]]
	Track[V, K]
	TrackEach[V, K]
	Checkable[K]
	Access[K, V]
	MapTransformable[K, V, map[K]V]
	Keys() Collection[K, []K, Iterator[K]]
	Values() Collection[V, []V, Iterator[V]]
}

//Iterator is the interface that provides iterate over elements of a collection.
type Iterator[T any] interface {
	//retrieves a next element and true or zero value of T and false if no more elements.
	Next() (T, bool)
	//returns an estimated internal storage capacity
	Cap() int
}

//PrevIterator is the Iterator that provides reverse iteration over elements of a collection.
type PrevIterator[T any] interface {
	Iterator[T]
	//retrieves a prev element and true or zero value of T and false if no more elements.
	Prev() (T, bool)
}

//DelIterator is the Iterator provides deleting of current element.
type DelIterator[T any] interface {
	Iterator[T]
	Delete() bool
}

//KVIterator is the interface that provides iterate over all key/value pair of a map.
type KVIterator[K, V any] interface {
	//retrieves next elements or zero values if no more elements
	Next() (K, V, bool)
	Cap() int
}

//Iterable is an iterator supplier interface
type Iterable[IT any] interface {
	Begin() IT
}

//Walk is the interface of a collection that provides traversing of the elements.
type Walk[IT any] interface {
	For(func(element IT) error) error
}

//WalkEach is the interface of a collection that provides traversing of the elements without error checking.
type WalkEach[T any] interface {
	ForEach(func(element T))
}

//Track is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.).
type Track[T any, P any] interface {
	Track(func(position P, element T) error) error
}

//TrackEach is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.) without error checking
type TrackEach[T any, P any] interface {
	TrackEach(func(position P, element T))
}

//Checkable is container with ability to check if an element is present.
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

//Addable is the interface that provides appending the collection by elements.
type Addable[T any] interface {
	Add(...T) bool
}

//Settable is the interface that provides replacing an element by its pointer (index or key).
type Settable[P any, V any] interface {
	Set(key P, value V) bool
}

//Deleteable is the interface that provides removing any elements from the collection.
type Deleteable[k any] interface {
	Delete(...k) bool
}

//Removable is the interface that provides removing an element by its pointer (index or key).
type Removable[P any, V any] interface {
	Remove(P) (V, bool)
}
