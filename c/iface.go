// Package c provides common types of containers, utility types and functions
package c

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = errors.New("Break")

// KeyVal provides extracing of a keys or values collection from key/value pairs
type KeyVal[Keys any, Vals any] interface {
	Keys() Keys
	Values() Vals
}

// Collection is the base interface of non-associative collections
type Collection[T any] interface {
	Iterable[T]
	ForLoop[T]
	ForEachLoop[T]
	SliceFactory[T]

	Reduce(merger func(T, T) T) T
}

// Filterable provides filtering content functionality
type Filterable[T, Stream, StreamBreakable any] interface {
	Filter(predicate func(T) bool) Stream
	Filt(predicate func(T) (bool, error)) StreamBreakable
}

// Convertable provides converaton of collection elements functionality
type Convertable[T, Stream, StreamBreakable any] interface {
	Convert(converter func(T) T) Stream
	Conv(converter func(T) (T, error)) StreamBreakable
}

// SliceFactory collects the elements of the collection into a slice
type SliceFactory[T any] interface {
	Slice() []T
	Append([]T) []T
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
	Next() (out T, ok bool)

	ForLoop[T]
	ForEachLoop[T]
	All(consumer func(T) bool)
}

// IterFor extends an iterator type by a 'Start' function implementation
type IterFor[T any, I Iterator[T]] interface {
	// Start is used with for loop construct.
	// Returns the iterator itself, the first element, and ok == false if the iteration must be completed.
	//
	// 	var i IterFor = ...
	// 	for i, val, ok := i.Start(); ok; val, ok = i.Next() {
	//  	_ = val
	//	}
	Start() (iterator I, val T, ok bool)
}

// Sized - storage interface with measurable capacity
type Sized interface {
	// returns an estimated internal storage capacity or -1 if the capacity cannot be calculated
	Cap() int
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

// Iterable is an iterator supplier interface
type Iterable[T any] interface {
	Iter() Iterator[T]
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
type TrackLoop[P any, T any] interface {
	// return ErrBreak for loop breaking
	Track(func(position P, element T) error) error
}

// TrackEachLoop is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.) without error checking
type TrackEachLoop[P any, T any] interface {
	TrackEach(func(position P, element T))
}

// Checkable is container with ability to check if an element is present.
type Checkable[T any] interface {
	Contains(T) bool
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
type AddableAll[Iterable any] interface {
	AddAll(Iterable)
}

// AddableAllNew provides appending the collection by elements retrieved from another collection
type AddableAllNew[Iterable any] interface {
	AddAllNew(Iterable) bool
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
type SettableMap[Map any] interface {
	SetMap(m Map)
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
type ImmutableMapConvert[M any] interface {
	Immutable() M
}

// Removable provides removing an element by its pointer (index or key).
type Removable[P any, V any] interface {
	Remove(P) (V, bool)
}

// Summable is a type that supports the operator +
type Summable interface {
	constraints.Ordered | constraints.Complex | string
}

// Number is a type that supports the operators +, -, /, *
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

type RangeFunc[T any] func(func(T) bool)
