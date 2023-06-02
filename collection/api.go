// Package collection consists of common operations of c.Iterable based collections
package collection

import (
	"golang.org/x/exp/constraints"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	kvstream "github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
)

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To any, I c.Iterable[From]](collection I, converter func(From) To) stream.Iter[To] {
	b := collection.Iter()
	return stream.New(loop.Convert(b.Next, converter).Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To any, I c.Iterable[From]](collection I, converter func(From) (To, error)) breakStream.Iter[To] {
	b := collection.Iter()
	return breakStream.New(breakLoop.Conv(breakLoop.From(b.Next), converter).Next)
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any, I c.Iterable[From]](collection I, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	b := collection.Iter()
	f := loop.FilterAndConvert(b.Next, filter, converter)
	return stream.New(f.Next)
}

// Flat returns a stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To any, I c.Iterable[From]](collection I, by func(From) []To) stream.Iter[To] {
	b := collection.Iter()
	f := loop.Flat(b.Next, by)
	return stream.New(f.Next)
}

// Flatt returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable, I c.Iterable[From]](collection I, flattener func(From) ([]To, error)) breakStream.Iter[To] {
	h := collection.Iter()
	f := breakLoop.Flatt(breakLoop.From(h.Next), flattener)
	return breakStream.New(f.Next)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any, I c.Iterable[From]](collection I, filter func(From) bool, flattener func(From) []To) stream.Iter[To] {
	b := collection.Iter()
	f := loop.FilterAndFlat(b.Next, filter, flattener)
	return stream.New(f.Next)
}

// Filter instantiates a stream that checks elements by the 'filter' function and returns successful ones
func Filter[T any, I c.Iterable[T]](collection I, filter func(T) bool) stream.Iter[T] {
	b := collection.Iter()
	f := loop.Filter(b.Next, filter)
	return stream.New(f.Next)
}

// NotNil instantiates a stream that filters nullable elements
func NotNil[T any, I c.Iterable[*T]](collection I) stream.Iter[*T] {
	return Filter(collection, not.Nil[T])
}

// PtrVal creates a stream that transform pointers to the values referenced referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any, I c.Iterable[*T]](collection I) stream.Iter[T] {
	return stream.New(loop.PtrVal(collection.Iter().Next).Next)
}

// NoNilPtrVal creates a stream that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any, I c.Iterable[*T]](collection I) stream.Iter[T] {
	return stream.New(loop.NoNilPtrVal(collection.Iter().Next).Next)
}

// NilSafe - convert.NilSafe filters not nil next, converts that ones, filters not nils after converting and returns them
func NilSafe[From, To any](next func() (*From, bool), converter func(*From) *To) stream.Iter[*To] {
	return stream.New(loop.ConvertCheck(next, func(f *From) (*To, bool) {
		if f != nil {
			if t := converter(f); t != nil {
				return t, true
			}
		}
		return nil, false
	}).Next)
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, I c.Iterable[T]](collection I, by func(T) K) kvstream.Iter[K, T, map[K][]T] {
	it := loop.NewKeyValuer(collection.Iter().Next, by, as.Is[T])
	return kvstream.New(it.Next, kvloop.Group[K, T])
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any, I c.Iterable[T]](collection I, predicate func(T) bool) (v T, ok bool) {
	i := collection.Iter()
	return loop.First(i.Next, predicate)
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any, I c.Iterable[T]](collection I, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	i := collection.Iter()
	return breakLoop.Firstt(breakLoop.From(i.Next), predicate)
}

// Sort sorts the specified sortable collection that contains orderable elements
func Sort[O any, S interface{ Sort(less slice.Less[T]) O }, T any, f constraints.Ordered](collection S, by func(T) f) O {
	return collection.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}
