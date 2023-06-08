package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/break/c"
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// ConvFiltIter is the array based Iterator thath provides converting of elements by a Converter with addition filtering of the elements by a Predicate.
type ConvFiltIter[From, To any] struct {
	array      unsafe.Pointer
	elemSize   uintptr
	size, i    int
	converter  func(From) (To, error)
	filterFrom func(From) (bool, error)
	filterTo   func(To) (bool, error)
}

var _ c.Iterator[any] = (*ConvFiltIter[any, any])(nil)

var _ c.IterFor[any, *ConvFiltIter[any, any]] = (*ConvFiltIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *ConvFiltIter[From, To]) For(walker func(element To) error) error {
	return loop.For(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *ConvFiltIter[From, To]) Next() (t To, ok bool, err error) {
	if i == nil || i.array == nil {
		return t, false, nil
	}
	next := func() (out From, ok bool, err error) {
		return nextFilt(i.array, i.size, i.elemSize, i.filterFrom, &i.i)
	}
	for {
		if v, ok, err := next(); err != nil || !ok {
			return t, false, err
		} else if t, err = i.converter(v); err != nil {
			return t, false, err
		} else if ok, err := i.filterTo(t); err != nil {
			return t, false, err
		} else if ok {
			return t, true, nil
		}
	}
}

// Cap returns the iterator capacity
func (i *ConvFiltIter[From, To]) Cap() int {
	return i.size
}

// Start is used with for loop construct like 'for i, val, ok, err := i.Start(); ok || err != nil ; val, ok, err = i.Next() { if err != nil { return err }}'
func (i *ConvFiltIter[From, To]) Start() (*ConvFiltIter[From, To], To, bool, error) {
	return startBreakIt[To](i)
}

// ConvIter is the array based Iterator thath provides converting of elements by a ConvIter.
type ConvIter[From, To any] struct {
	array     unsafe.Pointer
	elemSize  uintptr
	size, i   int
	converter func(From) (To, error)
}

var _ c.Iterator[any] = (*ConvIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *ConvIter[From, To]) For(walker func(element To) error) error {
	return loop.For(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *ConvIter[From, To]) Next() (t To, ok bool, err error) {
	if i.i < i.size {
		v := *(*From)(notsafe.GetArrayElemRef(i.array, i.i, i.elemSize))
		i.i++
		t, err = i.converter(v)
		return t, err == nil, err
	}
	return t, false, nil
}

// Cap returns the iterator capacity
func (i *ConvIter[From, To]) Cap() int {
	return i.size
}

// Start is used with for loop construct like 'for i, val, ok, err := i.Start(); ok || err != nil ; val, ok, err = i.Next() { if err != nil { return err }}'
func (i *ConvIter[From, To]) Start() (*ConvIter[From, To], To, bool, error) {
	return startBreakIt[To](i)
}

func nextFilt[T any](array unsafe.Pointer, size int, elemSize uintptr, filter func(T) (bool, error), index *int) (v T, ok bool, err error) {
	i := *index
	for ; i < size; i++ {
		v = *(*T)(notsafe.GetArrayElemRef(array, i, elemSize))
		if ok, err = filter(v); err != nil || ok {
			*index = i + 1
			return v, ok, err
		}
	}
	*index = i + 1
	return v, false, nil
}
