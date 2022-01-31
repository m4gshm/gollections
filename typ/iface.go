//Package typ provides common types of containers, utility types and functions.
package typ

import (
	"github.com/m4gshm/gollections/op"
)

// Container is the base interface of implementations.
// Where:
// 		C - a type of slice, map or chain
// 		IT - the Iterator or the KVIterator
type Container[C any, IT Iter] interface {
	Collect() C
	Iterable[IT]
}

//Collection is the interface of a finite-size container.
type Collection[T any, C any, IT Iter] interface {
	Container[C, IT]
	Walk[T]
	WalkEach[T]
}

//Vector is the interface of a container stores ordered elements, provides index access.
type Vector[T any] interface {
	Collection[T, []T, Iterator[T]]
	Track[T, int]
	TrackEach[T, int]
	Access[int, T]
	Transformable[T, []T, Iterator[T]]
	Len() int
}

//Set is the interface of a container provides uniqueness (does't insert duplicated values).
type Set[T any] interface {
	Collection[T, []T, Iterator[T]]
	Transformable[T, []T, Iterator[T]]
	Checkable[T]
	Len() int
}

//Map is the interface of a container provides access to elements by key.
type Map[k comparable, v any] interface {
	Collection[*KV[k, v], map[k]v, KVIterator[k, v]]
	Track[v, k]
	TrackEach[v, k]
	Checkable[k]
	Access[k, v]
	MapTransformable[k, v, map[k]v]
	Keys() Collection[k, []k, Iterator[k]]
	Values() Collection[v, []v, Iterator[v]]
	Len() int
}

//Iter is the base for Iterator and KVIterator.
type Iter interface {
	//checks ability on next element or error.
	HasNext() bool
}

//Iterator is the interface provides iterate over elements of a collection.
type Iterator[T any] interface {
	Iter
	//retrieves next element
	//must be called only after HasNext
	//may raise panic. Calls Get() to prevent panic and checks an iteration error.
	Next() T
	//retrieves next element or error
	//must be called only after HasNext
	Get() (T, error)
}

//KVIterator is the interface provides iterate over all key/value pair of a map.
type KVIterator[k, v any] interface {
	Iter
	//retrieves next element
	//must be called only after HasNext
	//may raise panic. Calls Get() to prevent panic and checks an iteration error.
	Next() (k, v)
	//retrieves next element or error
	//must be called only after HasNext
	Get() (k, v, error)
}

//Resetable is the interface of an object with resettable state (e.g. slice based iterator).
type Resetable interface {
	Reset()
}

//Iterable is an iterator supplier interface
type Iterable[IT Iter] interface {
	Begin() IT
}

//Walk touches elements of the collection.
type Walk[T any] interface {
	For(func(element T) error) error
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
type Transformable[T any, C any, IT Iterator[T]] interface {
	Filter(Predicate[T]) Pipe[T, C, IT]
	Map(Converter[T, T]) Pipe[T, C, IT]
}

//Pipe extends Transformable by finalize methods like ForEach, Collect or Reduce.
type Pipe[T any, C any, IT Iterator[T]] interface {
	Transformable[T, C, IT]
	Container[C, IT]
	Walk[T]
	WalkEach[T]
	Reduce(op.Binary[T]) T
}

//MapTransformable is the interface that provides limited kit of map transformation methods.
//The full kit of transformer functions are in the package 'c/map_'
type MapTransformable[k comparable, v any, m any] interface {
	Filter(BiPredicate[k, v]) MapPipe[k, v, m]
	Map(BiConverter[k, v, k, v]) MapPipe[k, v, m]

	FilterKey(Predicate[k]) MapPipe[k, v, m]
	MapKey(Converter[k, k]) MapPipe[k, v, m]

	FilterValue(Predicate[v]) MapPipe[k, v, m]
	MapValue(Converter[v, v]) MapPipe[k, v, m]
}

//MapPipe extends MapTransformable by finalize methods like ForEach, Collect or Reduce.
type MapPipe[k comparable, v any, m any] interface {
	MapTransformable[k, v, m]
	Container[m, KVIterator[k, v]]
}

//Access is the interface that provides access to an element by its pointer (index, key, coordinate, etc.)
type Access[K any, V any] interface {
	Get(K) (V, bool)
}
