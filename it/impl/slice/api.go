package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/notsafe"
)

// Map instantiates Iterator that converts elements with a converter and returns them
func Map[From, To any, FS ~[]From](elements FS, by c.Converter[From, To]) Convert[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return Convert[From, To]{array: array, size: size, elemSize: elemSize, by: by}
}

// MapFit additionally filters 'From' elements.
func MapFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], by c.Converter[From, To]) ConvertFit[From, To] {
	var (
		header   = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array    = unsafe.Pointer(header.Data)
		size     = header.Len
		elemSize = notsafe.GetTypeSize[From]()
	)
	return ConvertFit[From, To]{array: array, size: size, elemSize: elemSize, by: by, fit: fit}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To]) Flatten[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return Flatten[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: by}
}

// FlattFit additionally filters –'From' elements.
func FlattFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], flatt c.Flatter[From, To]) FlattenFit[From, To] {
	var (
		header       = notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elements))
		array        = unsafe.Pointer(header.Data)
		size         = header.Len
		elemSizeFrom = notsafe.GetTypeSize[From]()
		elemSizeTo   = notsafe.GetTypeSize[To]()
	)
	return FlattenFit[From, To]{arrayFrom: array, sizeFrom: size, elemSizeFrom: elemSizeFrom, elemSizeTo: elemSizeTo, flatt: flatt, fit: fit}
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[T any, TS ~[]T](elements TS, filter c.Predicate[T]) Fit[T] {
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
