// Package iter provides implementations of slice based itarators
package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/predicate/always"
	"github.com/m4gshm/gollections/slice/iter"
)

// Conv instantiates Iterator that converts elements with a converter and returns them
func Conv[FS ~[]From, From, To any](elements FS, converter func(From) (To, error)) *ConvertIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return &ConvertIter[From, To]{array: array, size: size, elemSize: elemSize, converter: converter}
}

// FiltAndConv additionally filters 'From' elements.
func FiltAndConv[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), converter func(From) (To, error)) *ConvFitIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return &ConvFitIter[From, To]{array: array, size: size, elemSize: elemSize, converter: converter, filterFrom: filter, filterTo: as.ErrTail(always.True[To])}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flat[FS ~[]From, From, To any](elements FS, flatter func(From) ([]To, error)) *Flatten[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return &Flatten[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: flatter}
}

// FilterAndFlatt additionally filters –'From' elements.
func FiltAndFlat[FS ~[]From, From, To any](elements FS, filter func(From) (bool, error), flatt func(From) ([]To, error)) *FlatFiltIter[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return &FlatFiltIter[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: flatt, filter: filter}
}

// Filt instantiates Iterator that checks elements by filters and returns successful ones.
func Filt[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) *FiltIter[T] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[T]()
	)
	return &FiltIter[T]{array: array, size: size, elemSize: elemSize, filter: filter}
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) *FiltIter[*T] {
	return Filt(elements, as.ErrTail(check.NotNil[T]))
}

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[TS ~[]T, T any, K, V any](elements TS, keyExtractor func(T) (K, error), valsExtractor func(T) (V, error)) *KeyValuer[T, K, V] {
	return &KeyValuer[T, K, V]{iter: iter.NewHead(elements), keyExtractor: keyExtractor, valExtractor: valsExtractor}
}

// NewMultipleKeyValuer creates instance of the MultipleKeyValuer
func NewMultipleKeyValuer[TS ~[]T, T any, K, V any](elements TS, keysExtractor func(T) ([]K, error), valsExtractor func(T) ([]V, error)) *MultipleKeyValuer[T, K, V] {
	return &MultipleKeyValuer[T, K, V]{iter: iter.NewHead(elements), keysExtractor: keysExtractor, valsExtractor: valsExtractor}
}

// ToKV transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func ToKV[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) *KeyValuer[T, K, V] {
	kv := NewKeyValuer(elements, keyExtractor, valExtractor)
	return kv
}

// ToKVs transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ToKVs[TS ~[]T, T, K, V any](elements TS, keysExtractor func(T) ([]K, error), valsExtractor func(T) ([]V, error)) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(elements, keysExtractor, valsExtractor)
	return kv
}

// FlattKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattKeys[TS ~[]T, T, K any](elements TS, keysExtractor func(T) ([]K, error)) *MultipleKeyValuer[T, K, T] {
	kv := NewMultipleKeyValuer(elements, keysExtractor, as.ErrTail(as.Slice[T]))
	return kv
}

// FlattValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattValues[TS ~[]T, T, V any](elements TS, valsExtractor func(T) ([]V, error)) *MultipleKeyValuer[T, T, V] {
	kv := NewMultipleKeyValuer(elements, as.ErrTail(as.Slice[T]), valsExtractor)
	return kv
}
