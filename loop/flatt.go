package loop

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
)

// FlatFilterIter is the Iterator wrapper that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlatFilterIter[From, To any] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	next          func() (From, bool)
	flattener     func(From) []To
	filterFrom    func(From) bool
	filterTo      func(To) bool
}

var _ c.Iterator[any] = (*FlatFilterIter[any, any])(nil)
var _ c.IterFor[any, *FlatFilterIter[any, any]] = (*FlatFilterIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *FlatFilterIter[From, To]) For(walker func(element To) error) error {
	return For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *FlatFilterIter[From, To]) ForEach(walker func(element To)) {
	ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *FlatFilterIter[From, To]) Next() (t To, ok bool) {
	if i == nil {
		return t, ok
	}

	next := i.next
	if next == nil {
		return t, false
	}

	for {
		if sizeTo := i.sizeTo; sizeTo > 0 {
			if indTo := i.indTo; indTo < sizeTo {
				i.indTo++
				t = *(*To)(notsafe.GetArrayElemRef(i.arrayTo, indTo, i.elemSizeTo))
				if ok := i.filterTo(t); ok {
					return t, true
				}
			}
			i.indTo = 0
			i.arrayTo = nil
			i.sizeTo = 0
		}

		if v, ok := next(); !ok {
			return t, false
		} else if i.filterFrom(v) {
			if elementsTo := i.flattener(v); len(elementsTo) > 0 {
				i.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				i.arrayTo = unsafe.Pointer(header.Data)
				i.sizeTo = header.Len
				t = *(*To)(notsafe.GetArrayElemRef(i.arrayTo, 0, i.elemSizeTo))
				if ok := i.filterTo(t); ok {
					return t, true
				}
			}
		}
	}
}

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (i *FlatFilterIter[From, To]) Start() (*FlatFilterIter[From, To], To, bool) {
	return startIt[To](i)
}

// FlatIter is the Iterator wrapper that converts an element to a slice and iterates over the elements of that slice.
// For example, FlatIter can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type FlatIter[From, To any] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	next          func() (From, bool)
	flattener     func(From) []To
}

var _ c.Iterator[any] = (*FlatIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *FlatIter[From, To]) For(walker func(element To) error) error {
	return For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *FlatIter[From, To]) ForEach(walker func(element To)) {
	ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *FlatIter[From, To]) Next() (To, bool) {
	sizeTo := i.sizeTo
	if sizeTo > 0 {
		if indTo := i.indTo; indTo < sizeTo {
			i.indTo++
			return *(*To)(notsafe.GetArrayElemRef(i.arrayTo, indTo, i.elemSizeTo)), true
		}
		i.indTo = 0
		i.arrayTo = nil
		i.sizeTo = 0
	}

	for {
		if v, ok := i.next(); !ok {
			var no To
			return no, false
		} else if elementsTo := i.flattener(v); len(elementsTo) > 0 {
			i.indTo = 1
			header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			i.arrayTo = unsafe.Pointer(header.Data)
			i.sizeTo = header.Len
			return *(*To)(notsafe.GetArrayElemRef(i.arrayTo, 0, i.elemSizeTo)), true
		}
	}
}

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (i *FlatIter[From, To]) Start() (*FlatIter[From, To], To, bool) {
	return startIt[To](i)
}
