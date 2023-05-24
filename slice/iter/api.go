// Package iter provides implementations of slice based itarators
package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/predicate/always"
	"github.com/m4gshm/gollections/slice"
)

// New is the main slice-based iterator constructor
func New[TS ~[]T, T any](elements TS) *slice.Iter[T] {
	return slice.NewIter(elements)
}

// Convert instantiates an iterator that converts elements with a converter and returns them
func Convert[FS ~[]From, From, To any](elements FS, converter func(From) To) *ConvertIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return &ConvertIter[From, To]{array: array, size: size, elemSize: elemSize, converter: converter}
}

// Conv instantiates an iterator that converts elements with a converter and returns them
func Conv[FS ~[]From, From, To any](elements FS, converter func(From) (To, error)) *ConvIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return &ConvIter[From, To]{array: array, size: size, elemSize: elemSize, converter: converter}
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) *ConvertFiltIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return &ConvertFiltIter[From, To]{array: array, size: size, elemSize: elemSize, converter: converter, filterFrom: filter, filterTo: always.True[To]}
}

// FiltAndConv additionally filters 'From' elements.
func FiltAndConv[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), converter func(From) (To, error)) *ConvFiltIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return &ConvFiltIter[From, To]{array: array, size: size, elemSize: elemSize, converter: converter, filterFrom: filter, filterTo: as.ErrTail(always.True[To])}
}

// Flat instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flat[FS ~[]From, From, To any](elements FS, by func(From) []To) *FlattIter[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return &FlattIter[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flattener: by}
}

// Flatt instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[FS ~[]From, From, To any](elements FS, flatter func(From) ([]To, error)) *FlatIter[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return &FlatIter[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flattener: flatter}
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[FS ~[]From, From, To any](elements FS, filter func(From) bool, flattener func(From) []To) *FlattenFiltIter[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return &FlattenFiltIter[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flattener: flattener, filter: filter}
}

// FiltAndFlat instantiates an iterator that filters elements by the 'filter' function, flattens elements and returns them.
func FiltAndFlat[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), flattener func(From) ([]To, error)) *FlatFiltIter[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return &FlatFiltIter[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flattener: flattener, filter: filter}
}

// Filter instantiates an iterator that checks elements by the 'filter' function and returns successful ones
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) *FilterIter[T] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[T]()
	)
	return &FilterIter[T]{array: array, size: size, elemSize: elemSize, filter: filter}
}

// Filt instantiates an iterator that checks elements by the 'filter' function and returns successful ones
func Filt[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) *FiltIter[T] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[T]()
	)
	return &FiltIter[T]{array: array, size: size, elemSize: elemSize, filter: filter}
}

// NotNil instantiates an iterator that filters nullable elements
func NotNil[T any, TRS ~[]*T](elements TRS) *FilterIter[*T] {
	return Filter(elements, not.Nil[T])
}

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[TS ~[]T, T any, K, V any](elements TS, keyExtractor func(T) K, valsExtractor func(T) V) *KeyValuerIter[T, K, V] {
	return &KeyValuerIter[T, K, V]{iter: slice.NewHead(elements), keyExtractor: keyExtractor, valExtractor: valsExtractor}
}

// NewMultipleKeyValuer creates instance of the MultipleKeyValuer
func NewMultipleKeyValuer[TS ~[]T, T any, K, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) *MultipleKeyValuer[T, K, V] {
	return &MultipleKeyValuer[T, K, V]{iter: slice.NewHead(elements), keysExtractor: keysExtractor, valsExtractor: valsExtractor}
}

// ToKV transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func ToKV[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) *KeyValuerIter[T, K, V] {
	kv := NewKeyValuer(elements, keyExtractor, valExtractor)
	return kv
}

// ToKVs transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ToKVs[TS ~[]T, T, K, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(elements, keysExtractor, valsExtractor)
	return kv
}

// FlattKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattKeys[TS ~[]T, T, K any](elements TS, keysExtractor func(T) []K) *MultipleKeyValuer[T, K, T] {
	kv := NewMultipleKeyValuer(elements, keysExtractor, convert.AsSlice[T])
	return kv
}

// FlattValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattValues[TS ~[]T, T, V any](elements TS, valsExtractor func(T) []V) *MultipleKeyValuer[T, T, V] {
	kv := NewMultipleKeyValuer(elements, convert.AsSlice[T], valsExtractor)
	return kv
}

// NewKeyVal creates instance of the KeyValuer
func NewKeyVal[TS ~[]T, T any, K, V any](elements TS, keyExtractor func(T) (K, error), valsExtractor func(T) (V, error)) *KeyValIter[T, K, V] {
	return &KeyValIter[T, K, V]{iter: slice.NewHead(elements), keyExtractor: keyExtractor, valExtractor: valsExtractor}
}
