// Package c provides common types of containers, utility types and functions
package c

import (
	"errors"

	"golang.org/x/exp/constraints"
)

// Break is the 'break' statement of the For, Track methods
var Break = errors.New("Break")

// Continue is an alias of the nil value used to continue iterating by For, Track methods.
var Continue error = nil

// Range provides the `Range` function used for iterating over a sequence of elements by `for e := range collection.Range`.
type Range[T any] interface {
	All(consumer func(T) bool)
}

// OrderedRange provides the `All` function used for iterating over an ordered sequence of elements by `for i, e := range collection.IAll`.
type OrderedRange[T any] interface {
	IAll(consumer func(int, T) bool)
}

// KVRange provides the `All` function used for iterating over a sequence of key\value pairs by `for k, v := range collection.All`.
type KVRange[K, V any] interface {
	All(consumer func(K, V) bool)
}

// Iterable is a loop supplier interface
type Iterable[T any, Loop ~func() (T, bool)] interface {
	Loop() Loop
}

// KeyVal provides access to all keys and values of a key/value based collection.
type KeyVal[K, V any] interface {
	Keys[K]
	Values[V]
}

// Keys provides access to all keys of a key/value based collection.
type Keys[K any] interface {
	Keys() K
}

// Values provides access to all values of a key/value based collection.
type Values[V any] interface {
	Values() V
}

// Collection is the base interface of non-associative collections
type Collection[T any] interface {
	Range[T]
	For[T]
	ForEach[T]
	SliceFactory[T]

	Reduce(merge func(T, T) T) T
	HasAny(func(T) bool) bool
}

// Filterable provides filtering content functionality
type Filterable[T any, Loop ~func() (T, bool), LoopErr ~func() (T, bool, error)] interface {
	Filter(predicate func(T) bool) Loop
	Filt(predicate func(T) (bool, error)) LoopErr
}

// Convertable provides converaton of collection elements functionality
type Convertable[T any, Loop ~func() (T, bool), LoopErr ~func() (T, bool, error)] interface {
	Convert(converter func(T) T) Loop
	Conv(converter func(T) (T, error)) LoopErr
}

// SliceFactory collects the elements of the collection into a slice
type SliceFactory[T any] interface {
	Slice() []T
	Append([]T) []T
}

// MapFactory collects the key/value pairs of the collection into a new map
type MapFactory[K comparable, V any, Map map[K]V | map[K][]V] interface {
	Map() Map
}

// Iterator provides iterate over elements of a collection
type Iterator[T any] interface {
	// Next returns the next element.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (out T, ok bool)

	For[T]
	ForEach[T]
}

// Sized - storage interface with measurable size
type Sized interface {
	// returns an estimated internal storage size or -1 if the size cannot be calculated
	Size() int
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

// For is the interface of a collection that provides traversing of the elements.
type For[IT any] interface {
	//For takes elements of the collection. Can be interrupt by returning Break.
	For(func(element IT) error) error
}

// ForEach is the interface of a collection that provides traversing of the elements without error checking.
type ForEach[T any] interface {
	// ForEach takes all elements of the collection
	ForEach(func(element T))
}

// Track is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.).
type Track[P any, T any] interface {
	// return Break for loop breaking
	Track(func(position P, element T) error) error
}

// TrackEach is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.) without error checking
type TrackEach[P any, T any] interface {
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
