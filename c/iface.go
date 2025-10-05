// Package c provides common types of containers, utility types and functions
package c

// Range provides an All function used for iterating over a sequence of elements by `for e := range collection.All`.
type Range[T any] interface {
	All(yield func(T) bool)
}

// OrderedRange provides an IAll function used for iterating over an ordered sequence of elements by `for i, e := range collection.IAll`.
type OrderedRange[T any] interface {
	IAll(yield func(int, T) bool)
}

// KVRange provides an All function used for iterating over a sequence of key\value pairs by `for k, v := range collection.All`.
type KVRange[K, V any] interface {
	All(yield func(K, V) bool)
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
	Sized
	Range[T]
	ForEach[T]
	SliceFactory[T]

	Head() (T, bool)
	First(func(T) bool) (T, bool)
	Reduce(merge func(T, T) T) T

	HasAny(func(T) bool) bool
}

// Filterable provides filtering content functionality
type Filterable[T any, Seq ~func(yield func(T) bool), SeqE ~func(yield func(T, error) bool)] interface {
	Filter(predicate func(T) bool) Seq
	Filt(predicate func(T) (bool, error)) SeqE
}

// Convertable provides converaton of collection elements functionality
type Convertable[T any, Seq ~func(yield func(T) bool), SeqE ~func(yield func(T, error) bool)] interface {
	Convert(converter func(T) T) Seq
	Conv(converter func(T) (T, error)) SeqE
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

	ForEach[T]
	Range[T]
}

// Sized - storage interface with measurable size
type Sized interface {
	// returns an estimated internal storage size or -1 if the size cannot be calculated
	Len() int
}

// ForEach is the interface of a collection that provides traversing of the elements without error checking.
type ForEach[T any] interface {
	// ForEach takes all elements of the collection
	ForEach(func(element T))
}

// TrackEach is the interface of a collection that provides traversing of the elements with position tracking (index, key, coordinates, etc.) without error checking.
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
