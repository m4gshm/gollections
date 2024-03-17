// Package collection consists of common operations of c.Iterable based collections
package collection

import (
	"golang.org/x/exp/constraints"

	breakloop "github.com/m4gshm/gollections/break/loop"
	breakstream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/comparer"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	kvstream "github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	loopconvert "github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/stream"
)

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To any, I c.Iterator[From], IT c.Iterable[From, I]](collection IT, converter func(From) To) stream.Iter[To] {
	b := collection.Iter()
	return stream.New(loop.Convert(b.Next, converter).Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To any, I c.Iterator[From], IT c.Iterable[From, I]](collection IT, converter func(From) (To, error)) breakstream.Iter[To] {
	b := collection.Iter()
	return breakstream.New(breakloop.Conv(breakloop.From(b.Next), converter).Next)
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any, I c.Iterator[From], IT c.Iterable[From, I]](collection IT, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	b := collection.Iter()
	f := loop.FilterAndConvert(b.Next, filter, converter)
	return stream.New(f.Next)
}

// Flat returns a stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To any, I c.Iterator[From], IT c.Iterable[From, I]](collection IT, by func(From) []To) stream.Iter[To] {
	b := collection.Iter()
	f := loop.Flat(b.Next, by)
	return stream.New(f.Next)
}

// Flatt returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable, I c.Iterator[From], IT c.Iterable[From, I]](collection IT, flattener func(From) ([]To, error)) breakstream.Iter[To] {
	h := collection.Iter()
	f := breakloop.Flatt(breakloop.From(h.Next), flattener)
	return breakstream.New(f.Next)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any, I c.Iterator[From], IT c.Iterable[From, I]](collection IT, filter func(From) bool, flattener func(From) []To) stream.Iter[To] {
	b := collection.Iter()
	f := loop.FilterAndFlat(b.Next, filter, flattener)
	return stream.New(f.Next)
}

// Filter instantiates a stream that checks elements by the 'filter' function and returns successful ones
func Filter[T any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, filter func(T) bool) stream.Iter[T] {
	b := collection.Iter()
	f := loop.Filter(b.Next, filter)
	return stream.New(f.Next)
}

// NotNil instantiates a stream that filters nullable elements
func NotNil[T any, I c.Iterator[*T], IT c.Iterable[*T, I]](collection IT) stream.Iter[*T] {
	return Filter(collection, not.Nil[T])
}

// PtrVal creates a stream that transform pointers to the values referenced referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any, I c.Iterator[*T], IT c.Iterable[*T, I]](collection IT) stream.Iter[T] {
	return stream.New(loop.PtrVal(collection.Iter().Next).Next)
}

// NoNilPtrVal creates a stream that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any, I c.Iterator[*T], IT c.Iterable[*T, I]](collection IT) stream.Iter[T] {
	return stream.New(loop.NoNilPtrVal(collection.Iter().Next).Next)
}

// NilSafe creates a stream that filters not nil elements, converts that ones, filters not nils after converting and returns them
func NilSafe[From, To any, I c.Iterator[*From], IT c.Iterable[*From, I]](collection IT, converter func(*From) *To) stream.Iter[*To] {
	h := collection.Iter()
	return stream.New(loopconvert.NilSafe(h.Next, converter).Next)
}

// KeyValue transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValue[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keyExtractor func(T) K, valExtractor func(T) V) loop.KeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.KeyValue(h.Next, keyExtractor, valExtractor)
}

// KeyValuee transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValuee[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) breakloop.KeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.KeyValuee(h.Next, keyExtractor, valExtractor)
}

// KeysValues transforms iterable elements to key/value iterator based on applying multiple keys, values extractor to the elements
func KeysValues[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keysExtractor func(T) []K, valsExtractor func(T) []V) *loop.MultipleKeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.NewMultipleKeyValuer(h.Next, keysExtractor, valsExtractor)
}

// KeysValue transforms iterable elements to key/value iterator based on applying keys, value extractor to the elements
func KeysValue[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keysExtractor func(T) []K, valExtractor func(T) V) *loop.MultipleKeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.KeysValues(h.Next, keysExtractor, func(t T) []V { return convert.AsSlice(valExtractor(t)) })
}

// KeysValuee transforms iterable elements to key/value iterator based on applying keys, value extractor to the elements
func KeysValuee[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keysExtractor func(T) ([]K, error), valExtractor func(T) (V, error)) *breakloop.MultipleKeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.KeysValuee(h.Next, keysExtractor, valExtractor)
}

// KeyValues transforms iterable elements to key/value iterator based on applying key, values extractor to the elements
func KeyValues[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keyExtractor func(T) K, valsExtractor func(T) []V) *loop.MultipleKeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.KeysValues(h.Next, func(t T) []K { return convert.AsSlice(keyExtractor(t)) }, valsExtractor)
}

// KeyValuess transforms iterable elements to key/value iterator based on applying key, values extractor to the elements
func KeyValuess[T, K, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keyExtractor func(T) (K, error), valsExtractor func(T) ([]V, error)) *breakloop.MultipleKeyValuer[T, K, V] {
	h := collection.Iter()
	return loop.KeyValuess(h.Next, keyExtractor, valsExtractor)
}

// ExtraVals transforms iterable elements to key/value iterator based on applying values extractor to the elements
func ExtraVals[T, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, valsExtractor func(T) []V) *loop.MultipleKeyValuer[T, T, V] {
	h := collection.Iter()
	return loop.KeyValues(h.Next, as.Is[T], valsExtractor)
}

// ExtraValss transforms iterable elements to key/value iterator based on applying values extractor to the elements
func ExtraValss[T, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, valsExtractor func(T) ([]V, error)) *breakloop.MultipleKeyValuer[T, T, V] {
	h := collection.Iter()
	return loop.KeyValuess(h.Next, as.ErrTail(as.Is[T]), valsExtractor)
}

// ExtraKeys transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeys[T, K any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keysExtractor func(T) []K) *loop.MultipleKeyValuer[T, K, T] {
	h := collection.Iter()
	return loop.KeysValue(h.Next, keysExtractor, as.Is[T])
}

// ExtraKeyss transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeyss[T, K any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keyExtractor func(T) (K, error)) *breakloop.MultipleKeyValuer[T, K, T] {
	h := collection.Iter()
	return loop.KeyValuess(h.Next, keyExtractor, as.ErrTail(convert.AsSlice[T]))
}

// ExtraKey transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKey[T, K any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keysExtractor func(T) K) loop.KeyValuer[T, K, T] {
	h := collection.Iter()
	return loop.KeyValue(h.Next, keysExtractor, as.Is[T])
}

// ExtraKeyy transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeyy[T, K any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, keyExtractor func(T) (K, error)) breakloop.KeyValuer[T, K, T] {
	h := collection.Iter()
	return loop.ExtraKeyy(h.Next, keyExtractor)
}

// ExtraValue transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraValue[T, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, valueExtractor func(T) V) loop.KeyValuer[T, T, V] {
	h := collection.Iter()
	return loop.ExtraValue(h.Next, valueExtractor)
}

// ExtraValuee transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraValuee[T, V any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, valExtractor func(T) (V, error)) breakloop.KeyValuer[T, T, V] {
	h := collection.Iter()
	return loop.ExtraValuee(h.Next, valExtractor)
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, by func(T) K) kvstream.Iter[K, T, map[K][]T] {
	it := loop.NewKeyValuer(collection.Iter().Next, by, as.Is[T])
	return kvstream.New(it.Next, kvloop.Group[K, T])
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, predicate func(T) bool) (v T, ok bool) {
	i := collection.Iter()
	return loop.First(i.Next, predicate)
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any, I c.Iterator[T], IT c.Iterable[T, I]](collection IT, predicate func(T) (bool, error)) (v T, ok bool, err error) {
	i := collection.Iter()
	return breakloop.Firstt(breakloop.From(i.Next), predicate)
}

// Sort sorts the specified sortable collection that contains orderable elements
func Sort[SC any, Cmp ~func(T, T) int, C interface {
	Sort(Cmp) SC
}, T any, O constraints.Ordered](collection C, order func(T) O) SC {
	return collection.Sort(comparer.Of(order))
}
