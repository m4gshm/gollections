// Package iter provides implementations of slice based itarators
package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/predicate/always"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[FS ~[]From, From, To any](elements FS, by func(From) To) ConvertIter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return ConvertIter[From, To]{array: array, size: size, elemSize: elemSize, converter: by}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) ConvertFit[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return ConvertFit[From, To]{array: array, size: size, elemSize: elemSize, converter: converter, filterFrom: filter, filterTo: always.True[To]}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[FS ~[]From, From, To any](elements FS, by func(From) []To) Flatten[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return Flatten[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: by}
}

// FilterAndFlatt additionally filters –'From' elements.
func FilterAndFlatt[FS ~[]From, From, To any](elements FS, filter func(From) bool, flatt func(From) []To) FlattenFitIter[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return FlattenFitIter[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: flatt, filter: filter}
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) FitIter[T] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[T]()
	)
	return FitIter[T]{array: array, size: size, elemSize: elemSize, filter: filter}
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) FitIter[*T] {
	return Filter(elements, check.NotNil[T])
}

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[TS ~[]T, T any, K, V any](elements TS, keyProducer func(T) K, valsProducer func(T) V) KeyValuer[T, K, V] {
	return KeyValuer[T, K, V]{iter: NewHead(elements), keyProducer: keyProducer, valProducer: valsProducer}
}

// NewMultipleKeyValuer creates instance of the MultipleKeyValuer
func NewMultipleKeyValuer[TS ~[]T, T any, K, V any](elements TS, keysProducer func(T) []K, valsProducer func(T) []V) MultipleKeyValuer[T, K, V] {
	return MultipleKeyValuer[T, K, V]{iter: NewHead(elements), keysProducer: keysProducer, valsProducer: valsProducer}
}
