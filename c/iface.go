//Package c provides common types of containers, utility types and functions.
package c

import (
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
)

// Container is the base interface of implementations.
// Where:
// 		col - a collection type (may be slice, map or chain)
// 		it - a iterator type  (Iterator, KVIterator)
type Container[col any, it Iter] interface {
	Collect() col
	Iterable[it]
}

//Collection is the interface of a finite-size container.
// Where:
//		t - any arbitrary type
// 		col - a collection type (may be slice, map or chain)
// 		it - a iterator type  (Iterator, KVIterator)
type Collection[t any, col any, it Iter] interface {
	Container[col, it]
	Walk[t]
	WalkEach[t]
}

//Vector is the interface of a container stores ordered elements, provides index access.
type Vector[t any] interface {
	Collection[t, []t, Iterator[t]]
	Track[t, int]
	TrackEach[t, int]
	Access[int, t]
	Transformable[t, []t, Iterator[t]]
	Len() int
}

//Set is the interface of a container provides uniqueness (does't insert duplicated values).
type Set[t any] interface {
	Collection[t, []t, Iterator[t]]
	Transformable[t, []t, Iterator[t]]
	Checkable[t]
	Len() int
}

//Map is the interface of a container provides access to elements by key.
type Map[k comparable, v any] interface {
	Collection[*map_.KV[k, v], map[k]v, KVIterator[k, v]]
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
type Iterator[t any] interface {
	Iter
	//retrieves next element or zero value if no more elements
	//must be called only after HasNext
	Next() t
}

//KVIterator is the interface provides iterate over all key/value pair of a map.
type KVIterator[k, v any] interface {
	Iter
	//retrieves next elements or zero values if no more elements
	//must be called only after HasNext
	Next() (k, v)
}

//Resetable is the interface of an object with resettable state (e.g. slice based iterator).
type Resetable interface {
	Reset()
}

//Iterable is an iterator supplier interface
type Iterable[it Iter] interface {
	Begin() it
}

//Walk touches elements of the collection.
type Walk[it any] interface {
	For(func(element it) error) error
}

//WalkEach touches all elements of the collection without error checking
type WalkEach[t any] interface {
	ForEach(func(element t))
}

//Track traverses container elements with position tracking (index, key, coordinates, etc.)
type Track[t any, p any] interface {
	Track(func(position p, element t) error) error
}

//TrackEach traverses container elements with position tracking (index, key, coordinates, etc.) without error checking
type TrackEach[T any, t any] interface {
	TrackEach(func(position t, element T))
}

//Checkable container with ability to check if an element is present.
type Checkable[t any] interface {
	Contains(t) bool
}

//Transformable is the interface that provides limited kit of container transformation methods.
//The full kit of transformer functions are in the package 'c'
type Transformable[t any, col any, it Iterator[t]] interface {
	Filter(Predicate[t]) Pipe[t, col, it]
	Map(Converter[t, t]) Pipe[t, col, it]
}

//Pipe extends Transformable by finalize methods like ForEach, Collect or Reduce.
type Pipe[t any, col any, it Iterator[t]] interface {
	Transformable[t, col, it]
	Container[col, it]
	Walk[t]
	WalkEach[t]
	Reduce(op.Binary[t]) t
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
// Where:
//		p - a type of pointer to a value (index, map key, coordinates)
// 		a - any arbitrary type of the value
type Access[p any, v any] interface {
	Get(p) (v, bool)
}
