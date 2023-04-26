// Package c provides common types of containers, utility types and functions
package c

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = errors.New("Break")

// Vector - collection interface that provides elements order and access by index to the elements.
type Vector[T any] interface {
	Collection[T]

	TrackLoop[T, int]
	TrackEachLoop[T, int]

	Access[int, T]

	Len() int
	IsEmpty() bool
}

// Set - collection interface that ensures the uniqueness of elements (does not insert duplicate values).
type Set[T any] interface {
	Collection[T]
	Checkable[T]

	Len() int
	IsEmpty() bool
}

// Map - collection interface that stores key/value pairs and provide access to an element by its key
type Map[K comparable, V any] interface {
	KVCollection[K, V, map[K]V]
	Checkable[K]
	Access[K, V]

	Keys() Collection[K]
	Values() Collection[V]

	Len() int
	IsEmpty() bool
}

// Collection is the base interface of non-associative collections
type Collection[T any] interface {
	Iterable[T]
	ForLoop[T]
	ForEachLoop[T]
	SliceFactory[T]

	Filter(predicate func(T) bool) Stream[T]
	Convert(converter func(T) T) Stream[T]

	Reduce(merger func(T, T) T) T
	HasAny(predicate func(T) bool) bool
}

// KVCollection is the base interface of associative collections
type KVCollection[K comparable, V any, M map[K]V | map[K][]V] interface {
	TrackLoop[V, K]
	TrackEachLoop[V, K]
	KVIterable[K, V]
	KVTransformable[K, V, M]
	MapFactory[K, V, M]

	Reduce(merger func(K, V, K, V) (K, V)) (K, V)
	HasAny(predicate func(K, V) bool) bool
}

// SliceFactory collects the elements of the collection into a slice
type SliceFactory[T any] interface {
	Slice() []T
}

// MapFactory collects the key/value pairs of the collection into a map
type MapFactory[K comparable, V any, Map map[K]V | map[K][]V] interface {
	Map() Map
}

// Iterator provides iterate over elements of a collection
type Iterator[T any] interface {
	// Next returns the next element.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (T, bool)

	ForLoop[T]
	ForEachLoop[T]
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
	TrackLoop[V, K]
	TrackEachLoop[V, K]
}

// KVIteratorBreakable provides iterate over key/value pairs, where an iteration can be interrupted by an error
type KVIteratorBreakable[K, V any] interface {
	KVIterator[K, V]
	//returns an iteration abort error
	Error() error
}

// Iterable is an iterator supplier interface
type Iterable[T any] interface {
	Begin() Iterator[T]
}

// KVIterable is an iterator supplier interface
type KVIterable[K, V any] interface {
	Begin() KVIterator[K, V]
}

// ForLoop is the interface of a collection that provides traversing of the elements.
type ForLoop[IT any] interface {
	//For takes elements of the collection. Can be interrupt by returning ErrBreak.
	For(func(element IT) error) error
}

// ForEachLoop is the interface of a collection that provides traversing of the elements without error checking.
type ForEachLoop[T any] interface {
	// ForEach takes all elements of the collection
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

// Stream is collection or stream of elements in transformation state.
type Stream[T any] interface {
	Iterator[T]
	Collection[T]
}

// KVTransformable provides limited kit of map transformation methods.
// The full kit of transformer functions are in the package 'c/map_'
type KVTransformable[K comparable, V any, Map map[K]V | map[K][]V] interface {
	Filter(predicate func(K, V) bool) KVStream[K, V, Map]
	Convert(converter func(K, V) (K, V)) KVStream[K, V, Map]

	FilterKey(predicate func(K) bool) KVStream[K, V, Map]
	ConvertKey(converter func(K) K) KVStream[K, V, Map]

	FilterValue(predicate func(V) bool) KVStream[K, V, Map]
	ConvertValue(converter func(V) V) KVStream[K, V, Map]
}

// KVStream is map or key/value stream of elements in transformation state.
type KVStream[K comparable, V any, M map[K]V | map[K][]V] interface {
	KVIterator[K, V]
	KVCollection[K, V, M]
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
	AddOneNew(T) bool
}

// AddableAll provides appending the collection by elements retrieved from another collection
type AddableAll[T any] interface {
	AddAll(Iterable[T])
}

// AddableAllNew provides appending the collection by elements retrieved from another collection
type AddableAllNew[T any] interface {
	AddAllNew(Iterable[T]) bool
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
