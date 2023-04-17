// Package c provides common types of containers, utility types and functions.
package c

import (
	"github.com/m4gshm/gollections/loop"
	"golang.org/x/exp/constraints"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = loop.ErrBreak

// Container is the base interface of implementations.
// Where:
//
//	Collection - a collection type (may be slice, map or chain)
//	IT - a iterator type  (Iterator, KVIterator)
type Container[Collection any, IT any] interface {
	Collect() Collection
	Iterable[IT]
}

// Collection is the interface of a finite-size container.
// Where:
//
//	T - any arbitrary type
//	Collection - a collection type (may be slice, map or chain)
//	IT - a iterator type  (Iterator, KVIterator)
type Collection[T any, Collection any, IT any] interface {
	Container[Collection, IT]
	ForLoop[T]
	ForEachLoop[T]
	Len() int
	IsEmpty() bool
}

// Vector - collection interface that provides elements order and access by index to the elements.
type Vector[T any] interface {
	Collection[T, []T, Iterator[T]]
	TrackLoop[T, int]
	TrackEachLoop[T, int]
	Access[int, T]
	Transformable[T, []T]
}

// Set - collection interface that ensures the uniqueness of elements (does not insert duplicate values).
type Set[T any] interface {
	Collection[T, []T, Iterator[T]]
	Transformable[T, []T]
	Checkable[T]
}

// Map - collection interface that stores key/value pairs and provide access to an element by its key
type Map[K comparable, V any] interface {
	Collection[KV[K, V], map[K]V, KVIterator[K, V]]
	TrackLoop[V, K]
	TrackEachLoop[V, K]
	Checkable[K]
	Access[K, V]
	MapTransformable[K, V, map[K]V]
	Keys() Collection[K, []K, Iterator[K]]
	Values() Collection[V, []V, Iterator[V]]
}

// Iterator provides iterate over elements of a collection
type Iterator[T any] interface {
	// retrieves a next element and true or zero value of T and false if no more elements
	Next() (T, bool)
}

// Sized - storage interface with measurable capacity
type Sized interface {
	// returns an estimated internal storage capacity or -1 if the capacity cannot be calculated
	Cap() int
}

// IteratorBreakable provides iterate over elements of a source, where an iteration can be interrupted by an error.
type IteratorBreakable[T any] interface {
	Iterator[T]
	//returns an iteration abort error
	Error() error
}

// PrevIterator is the Iterator that provides reverse iteration over elements of a collection
type PrevIterator[T any] interface {
	Iterator[T]
	//retrieves a prev element and true or zero value of T and false if no more elements
	Prev() (T, bool)
}

// DelIterator is the Iterator provides deleting of current element.
type DelIterator[T any] interface {
	Iterator[T]
	Delete()
}

// KVIterator provides iterate over key/value pairs
type KVIterator[K, V any] interface {
	//retrieves next elements or zero values if no more elements
	Next() (K, V, bool)
}

// KVIteratorBreakable provides iterate over key/value pairs, where an iteration can be interrupted by an error
type KVIteratorBreakable[K, V any] interface {
	KVIterator[K, V]
	//returns an iteration abort error
	Error() error
}

// Iterable is an iterator supplier interface
type Iterable[IT any] interface {
	Begin() IT
}

// ForLoop is the interface of a collection that provides traversing of the elements.
type ForLoop[IT any] interface {
	// return ErrBreak for loop breaking
	For(func(element IT) error) error
}

// ForEachLoop is the interface of a collection that provides traversing of the elements without error checking.
type ForEachLoop[T any] interface {
	ForEach(func(element T))
}

// TrackLoop is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.).
type TrackLoop[T any, P any] interface {
	// return ErrBreak for loop breaking
	Track(func(position P, element T) error) error
}

// TrackEachLoop is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.) without error checking
type TrackEachLoop[T any, P any] interface {
	TrackEach(func(position P, element T))
}

// Checkable is container with ability to check if an element is present.
type Checkable[T any] interface {
	Contains(T) bool
}

// Transformable provides limited kit of container transformation methods.
// The full kit of transformer functions are in the package 'c'
type Transformable[T any, Collection any] interface {
	Filter(func(T) bool) Pipe[T, Collection]
	Convert(func(T) T) Pipe[T, Collection]
}

// Pipe extends Transformable by finalize methods like ForEach, Collect or Reduce.
type Pipe[T any, Collection any] interface {
	Transformable[T, Collection]
	Container[Collection, Iterator[T]]
	ForLoop[T]
	ForEachLoop[T]
	Reduce(func(T, T) T) T
}

// MapTransformable provides limited kit of map transformation methods.
// The full kit of transformer functions are in the package 'c/map_'
type MapTransformable[K comparable, V any, Map any] interface {
	Filter(func(K, V) bool) MapPipe[K, V, Map]
	Convert(func(K, V) (K, V)) MapPipe[K, V, Map]

	FilterKey(func(K) bool) MapPipe[K, V, Map]
	ConvertKey(func(K) K) MapPipe[K, V, Map]

	FilterValue(func(V) bool) MapPipe[K, V, Map]
	ConvertValue(func(V) V) MapPipe[K, V, Map]
}

// MapPipe extends MapTransformable by finalize methods like ForEach, Collect or Reduce.
type MapPipe[K comparable, V any, Map any] interface {
	MapTransformable[K, V, Map]
	Container[Map, KVIterator[K, V]]
}

// Access provides access to an element by its pointer (index, key, coordinate, etc.)
// Where:
//
//	P - a type of pointer to a value (index, map key, coordinates)
//	V - any arbitrary type of the value
type Access[P any, V any] interface {
	Get(P) (V, bool)
}

// Addable provides appending the collection by elements.
type Addable[T any] interface {
	Add(...T)
	AddOne(T)
}

// AddableNew provides appending the collection by elements.
type AddableNew[T any] interface {
	AddNew(...T) bool
	AddNewOne(T) bool
}

// Settable provides element insertion or replacement by its pointer (index or key).
type Settable[P any, V any] interface {
	Set(key P, value V)
}

// SettableNew provides element insertion by its pointer (index or key) only if the specified place is not occupied.
type SettableNew[P any, V any] interface {
	SetNew(key P, value V) bool
}

// SettableMap provides element insertion or replacement with an equal key element of a map.
type SettableMap[K comparable, V any] interface {
	SetMap(m Map[K, V])
}

// Deleteable provides removing any elements from the collection.
type Deleteable[k any] interface {
	Delete(...k)
	DeleteOne(k)
}

// DeleteableVerify provides removing any elements from the collection.
type DeleteableVerify[k any] interface {
	DeleteActual(...k) bool
	DeleteActualOne(k) bool
}

// ImmutableMapConvert provides converting to an immutable map instance.
type ImmutableMapConvert[K comparable, V any, M Map[K, V]] interface {
	Immutable() M
}

// Removable provides removing an element by its pointer (index or key).
type Removable[P any, V any] interface {
	Remove(P) (V, bool)
}

// Summable is a type that supports the operator +
type Summable interface {
	constraints.Ordered | constraints.Complex
}

// Number is a type that supports the operators +, -, /, *
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Quaternary is an operation with four arguments
type Quaternary[t1, t2 any] func(t1, t2, t1, t2) (t1, t2)
