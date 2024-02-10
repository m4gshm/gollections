package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// FlattenFiltIter is the array based Iterator impelementation that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlattenFiltIter[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo, cap      int
	flattener                func(From) []To
	filter                   func(From) bool
}

var _ c.Iterator[any] = (*FlattenFiltIter[any, any])(nil)
var _ c.IterFor[any, *FlattenFiltIter[any, any]] = (*FlattenFiltIter[any, any])(nil)

func (f *FlattenFiltIter[From, To]) All(yield func(element To) bool) {
	loop.All(f.Next, yield)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *FlattenFiltIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f *FlattenFiltIter[From, To]) ForEach(walker func(element To)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f *FlattenFiltIter[From, To]) Next() (To, bool) {
	if sizeTo := f.sizeTo; sizeTo > 0 {
		if indTo := f.indTo; indTo < sizeTo {
			f.indTo++
			return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, indTo, f.elemSizeTo)), true
		}
		f.indTo = 0
		f.arrayTo = nil
		f.sizeTo = 0
	}

	for indFrom := f.indFrom; indFrom < f.sizeFrom; indFrom++ {
		if v := *(*From)(notsafe.GetArrayElemRef(f.arrayFrom, indFrom, f.elemSizeFrom)); f.filter(v) {
			if elementsTo := f.flattener(v); len(elementsTo) > 0 {
				f.indFrom = indFrom + 1
				f.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				f.arrayTo = unsafe.Pointer(header.Data)
				f.sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, 0, f.elemSizeTo)), true
			}
		}
	}
	var no To
	return no, false
}

// Cap returns the iterator capacity
func (f *FlattenFiltIter[From, To]) Cap() int {
	return f.cap
}

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (f *FlattenFiltIter[From, To]) Start() (*FlattenFiltIter[From, To], To, bool) {
	return startIt[To](f)
}

// FlattIter is the array based Iterator impelementation that converts an element to a slice and iterates over the elements of that slice.
// For example, FlattIter can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type FlattIter[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo           int
	flattener                func(From) []To
}

var _ c.Iterator[any] = (*FlattIter[any, any])(nil)
var _ c.IterFor[any, *FlattIter[any, any]] = (*FlattIter[any, any])(nil)

func (f *FlattIter[From, To]) All(yield func(element To) bool) {
	loop.All(f.Next, yield)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak.
func (f *FlattIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f *FlattIter[From, To]) ForEach(walker func(element To)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f *FlattIter[From, To]) Next() (To, bool) {
	sizeTo := f.sizeTo
	if sizeTo > 0 {
		if indTo := f.indTo; indTo < sizeTo {
			f.indTo++
			return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, indTo, f.elemSizeTo)), true
		}
		f.indTo = 0
		f.arrayTo = nil
		f.sizeTo = 0
	}
	for indFrom := f.indFrom; indFrom < f.sizeFrom; indFrom++ {
		if elementsTo := f.flattener(*(*From)(notsafe.GetArrayElemRef(f.arrayFrom, indFrom, f.elemSizeFrom))); len(elementsTo) > 0 {
			f.indFrom = indFrom + 1
			f.indTo = 1
			header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			f.arrayTo = unsafe.Pointer(header.Data)
			f.sizeTo = header.Len
			return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, 0, f.elemSizeTo)), true
		}
	}
	var no To
	return no, false
}

// Cap returns the iterator capacity
func (f *FlattIter[From, To]) Cap() int {
	return f.sizeFrom
}

// Start is used with for loop construct like 'for i, k, v, ok := i.Start(); ok; k, v, ok = i.Next() { }'
func (f *FlattIter[From, To]) Start() (*FlattIter[From, To], To, bool) {
	return startIt[To](f)
}
