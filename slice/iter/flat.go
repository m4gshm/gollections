package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/break/c"
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// FlatFiltIter is the array based Iterator impelementation that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlatFiltIter[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo, cap      int
	flatt                    func(From) ([]To, error)
	filter                   func(From) (bool, error)
}

var _ c.Iterator[any] = (*FlatFiltIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *FlatFiltIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f *FlatFiltIter[From, To]) Next() (t To, ok bool, err error) {
	if sizeTo := f.sizeTo; sizeTo > 0 {
		if indTo := f.indTo; indTo < sizeTo {
			f.indTo++
			return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, indTo, f.elemSizeTo)), true, nil
		}
		f.indTo = 0
		f.arrayTo = nil
		f.sizeTo = 0
	}

	for indFrom := f.indFrom; indFrom < f.sizeFrom; indFrom++ {
		v := *(*From)(notsafe.GetArrayElemRef(f.arrayFrom, indFrom, f.elemSizeFrom))
		if ok, err := f.filter(v); err != nil {
			return t, false, err
		} else if ok {
			if elementsTo, err := f.flatt(v); err != nil {
				return t, false, err
			} else if len(elementsTo) > 0 {
				f.indFrom = indFrom + 1
				f.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				f.arrayTo = unsafe.Pointer(header.Data)
				f.sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, 0, f.elemSizeTo)), true, nil
			}
		}
	}
	return t, false, nil
}

// Cap returns the iterator capacity
func (f *FlatFiltIter[From, To]) Cap() int {
	return f.cap
}

// FlatIter is the array based Iterator impelementation that converts an element to a slice and iterates over the elements of that slice.
// For example, FlatIter can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type FlatIter[From, To any] struct {
	arrayFrom, arrayTo       unsafe.Pointer
	elemSizeFrom, elemSizeTo uintptr
	sizeFrom, sizeTo         int
	indFrom, indTo           int
	flatt                    func(From) ([]To, error)
}

var _ c.Iterator[any] = (*FlatIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *FlatIter[From, To]) For(walker func(element To) error) error {
	return loop.For(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f *FlatIter[From, To]) Next() (t To, ok bool, err error) {
	sizeTo := f.sizeTo
	if sizeTo > 0 {
		if indTo := f.indTo; indTo < sizeTo {
			f.indTo++
			return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, indTo, f.elemSizeTo)), true, nil
		}
		f.indTo = 0
		f.arrayTo = nil
		f.sizeTo = 0
	}
	for indFrom := f.indFrom; indFrom < f.sizeFrom; indFrom++ {
		vf := *(*From)(notsafe.GetArrayElemRef(f.arrayFrom, indFrom, f.elemSizeFrom))
		if elementsTo, err := f.flatt(vf); err != nil {
			return t, false, err
		} else if len(elementsTo) > 0 {
			f.indFrom = indFrom + 1
			f.indTo = 1
			header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			f.arrayTo = unsafe.Pointer(header.Data)
			f.sizeTo = header.Len
			return *(*To)(notsafe.GetArrayElemRef(f.arrayTo, 0, f.elemSizeTo)), true, nil
		}
	}
	return t, false, nil
}

// Cap returns the iterator capacity
func (f *FlatIter[From, To]) Cap() int {
	return f.sizeFrom
}
