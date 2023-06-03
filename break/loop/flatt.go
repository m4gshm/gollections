package loop

import (
	"unsafe"

	"github.com/m4gshm/gollections/break/c"
	"github.com/m4gshm/gollections/notsafe"
)

// FlattFiltIter is the Iterator wrapper that converts an element to a slice with addition filtering of the element by a Predicate and iterates over the slice.
type FlattFiltIter[From, To any] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	next          func() (From, bool, error)
	flattener     func(From) ([]To, error)
	filterFrom    func(From) (bool, error)
	filterTo      func(To) (bool, error)
}

var _ c.Iterator[any] = (*FlattFiltIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *FlattFiltIter[From, To]) For(walker func(element To) error) error {
	return For(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *FlattFiltIter[From, To]) Next() (t To, ok bool, err error) {
	if i == nil {
		return t, false, nil
	}

	next := i.next
	if next == nil {
		return t, false, nil
	}

	for {
		if sizeTo := i.sizeTo; sizeTo > 0 {
			if indTo := i.indTo; indTo < sizeTo {
				i.indTo++
				t = *(*To)(notsafe.GetArrayElemRef(i.arrayTo, indTo, i.elemSizeTo))
				if ok, err := i.filterTo(t); err != nil {
					return t, false, err
				} else if ok {
					return t, true, nil
				}
			}
			i.indTo = 0
			i.arrayTo = nil
			i.sizeTo = 0
		}

		if v, ok, err := next(); err != nil || !ok {
			return t, false, err
		} else if ok, err := i.filterFrom(v); err != nil {
			return t, false, err
		} else if ok {
			if elementsTo, err := i.flattener(v); err != nil {
				return t, false, err
			} else if len(elementsTo) > 0 {
				i.indTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				i.arrayTo = unsafe.Pointer(header.Data)
				i.sizeTo = header.Len
				t = *(*To)(notsafe.GetArrayElemRef(i.arrayTo, 0, i.elemSizeTo))
				if ok, err := i.filterTo(t); err != nil || ok {
					return t, ok, err
				}
			}
		}
	}
}

func (i *FlattFiltIter[From, To]) Start() (*FlattFiltIter[From, To], To, bool, error) {
	return startIt[To](i)
}

// FlatIter is the Iterator wrapper that converts an element to a slice and iterates over the elements of that slice.
// For example, FlatIter can be used to iterate over all the elements of a multi-dimensional array as if it were a one-dimensional array ([][]int -> []int).
type FlatIter[From, To any] struct {
	arrayTo       unsafe.Pointer
	elemSizeTo    uintptr
	indTo, sizeTo int
	next          func() (From, bool, error)
	flattener     func(From) ([]To, error)
}

var _ c.Iterator[any] = (*FlatIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *FlatIter[From, To]) For(walker func(element To) error) error {
	return For(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *FlatIter[From, To]) Next() (t To, ok bool, err error) {
	sizeTo := i.sizeTo
	if sizeTo > 0 {
		if indTo := i.indTo; indTo < sizeTo {
			i.indTo++
			return *(*To)(notsafe.GetArrayElemRef(i.arrayTo, indTo, i.elemSizeTo)), true, nil
		}
		i.indTo = 0
		i.arrayTo = nil
		i.sizeTo = 0
	}

	for {
		if v, ok, err := i.next(); err != nil {
			return t, false, err
		} else if !ok {
			return t, false, nil
		} else if elementsTo, err := i.flattener(v); err != nil {
			return t, false, err
		} else if len(elementsTo) > 0 {
			i.indTo = 1
			header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
			i.arrayTo = unsafe.Pointer(header.Data)
			i.sizeTo = header.Len
			return *(*To)(notsafe.GetArrayElemRef(i.arrayTo, 0, i.elemSizeTo)), true, nil
		}
	}
}

func (i *FlatIter[From, To]) Start() (*FlatIter[From, To], To, bool, error) {
	return startIt[To](i)
}
