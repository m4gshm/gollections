package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/predicate"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[FS ~[]From, From, To any](elements FS, by c.Converter[From, To]) Converter[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return Converter[From, To]{array: array, size: size, elemSize: elemSize, by: by}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter predicate.Predicate[From], converter c.Converter[From, To]) ConvertFit[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return ConvertFit[From, To]{array: array, size: size, elemSize: elemSize, by: converter, filter: filter}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[FS ~[]From, From, To any](elements FS, by c.Flatter[From, To]) Flatten[From, To] {
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
func FilterAndFlatt[FS ~[]From, From, To any](elements FS, filter predicate.Predicate[From], flatt c.Flatter[From, To]) FlattenFit[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return FlattenFit[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: flatt, filter: filter}
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[TS ~[]T, T any](elements TS, filter predicate.Predicate[T]) Fit[T] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[T]()
	)
	return Fit[T]{array: array, size: size, elemSize: elemSize, by: filter}
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) Fit[*T] {
	return Filter(elements, check.NotNil[T])
}
