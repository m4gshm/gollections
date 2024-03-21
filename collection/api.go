// Package collection consists of common operations of c.Iterable based collections
package collection

import (
	"golang.org/x/exp/constraints"

	breakloop "github.com/m4gshm/gollections/break/loop"
	breakstream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/comparer"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	kvstream "github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	loopconvert "github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/stream"
)

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To any, IT c.Iterable[From]](collection IT, converter func(From) To) stream.Iter[To] {
	b := collection.Loop()
	return stream.New(loop.Convert(b, converter))
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To any, IT c.Iterable[From]](collection IT, converter func(From) (To, error)) breakstream.Iter[To] {
	b := collection.Loop()
	return breakstream.New(breakloop.Conv(breakloop.From(b), converter))
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any, IT c.Iterable[From]](collection IT, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	b := collection.Loop()
	f := loop.FilterAndConvert(b, filter, converter)
	return stream.New(f)
}

// Flat returns a stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To any, IT c.Iterable[From]](collection IT, by func(From) []To) stream.Iter[To] {
	b := collection.Loop()
	f := loop.Flat(b, by)
	return stream.New(f)
}

// Flatt returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable, IT c.Iterable[From]](collection IT, flattener func(From) ([]To, error)) breakstream.Iter[To] {
	h := collection.Loop()
	f := breakloop.Flatt(breakloop.From(h), flattener)
	return breakstream.New(f)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any, IT c.Iterable[From]](collection IT, filter func(From) bool, flattener func(From) []To) stream.Iter[To] {
	b := collection.Loop()
	f := loop.FilterAndFlat(b, filter, flattener)
	return stream.New(f)
}

// Filter instantiates a stream that checks elements by the 'filter' function and returns successful ones
func Filter[T any, IT c.Iterable[T]](collection IT, filter func(T) bool) stream.Iter[T] {
	b := collection.Loop()
	f := loop.Filter(b, filter)
	return stream.New(f)
}

// NotNil instantiates a stream that filters nullable elements
func NotNil[T any, IT c.Iterable[*T]](collection IT) stream.Iter[*T] {
	return Filter(collection, not.Nil[T])
}

// PtrVal creates a stream that transform pointers to the values referenced referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any, IT c.Iterable[*T]](collection IT) stream.Iter[T] {
	return stream.New(loop.PtrVal(collection.Loop()))
}

// NoNilPtrVal creates a stream that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any, IT c.Iterable[*T]](collection IT) stream.Iter[T] {
	return stream.New(loop.NoNilPtrVal(collection.Loop()))
}

// NilSafe creates a stream that filters not nil elements, converts that ones, filters not nils after converting and returns them
func NilSafe[From, To any, IT c.Iterable[*From]](collection IT, converter func(*From) *To) stream.Iter[*To] {
	h := collection.Loop()
	return stream.New(loopconvert.NilSafe(h, converter))
}

// KeyValue transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValue[T any, K comparable, V any, IT c.Iterable[T]](collection IT, keyExtractor func(T) K, valExtractor func(T) V) kvstream.Iter[K, V, map[K][]V] {
	h := collection.Loop()
	return kvstream.New(loop.KeyValue(h, keyExtractor, valExtractor), kvloop.Group)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any, IT c.Iterable[T]](collection IT, predicate func(T) bool) (v T, ok bool) {
	i := collection.Loop()
	return loop.First(i, predicate)
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any, IT c.Iterable[T]](collection IT, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	i := collection.Loop()
	return breakloop.Firstt(breakloop.From(i), predicate)
}

// Sort sorts the specified sortable collection that contains orderable elements
func Sort[SC any, Cmp ~func(T, T) int, C interface {
	Sort(Cmp) SC
}, T any, O constraints.Ordered](collection C, order func(T) O) SC {
	return collection.Sort(comparer.Of(order))
}
