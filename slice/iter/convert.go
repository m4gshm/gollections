package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
)

// ConvertFiltIter is the array based Iterator thath provides converting of elements by a Converter with addition filtering of the elements by a Predicate.
type ConvertFiltIter[From, To any] struct {
	array      unsafe.Pointer
	elemSize   uintptr
	size, i    int
	converter  func(From) To
	filterFrom func(From) bool
	filterTo   func(To) bool
}

var _ c.Iterator[any] = (*ConvertFiltIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *ConvertFiltIter[From, To]) For(walker func(element To) error) error {
	return loop.For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *ConvertFiltIter[From, To]) ForEach(walker func(element To)) {
	loop.ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *ConvertFiltIter[From, To]) Next() (t To, ok bool) {
	if i == nil || i.array == nil {
		return t, false
	}
	next := func() (From, bool) { return nextFiltered(i.array, i.size, i.elemSize, i.filterFrom, &i.i) }
	for v, ok := next(); ok; v, ok = next() {
		if t = i.converter(v); i.filterTo(t) {
			return t, true
		}
	}
	return t, false
}

// Cap returns the iterator capacity
func (i *ConvertFiltIter[From, To]) Cap() int {
	return i.size
}

func (i *ConvertFiltIter[From, To]) Start() (*ConvertFiltIter[From, To], To, bool) {
	return startIt[To](i)
}

// ConvertIter is the array based Iterator thath provides converting of elements by a ConvertIter.
type ConvertIter[From, To any] struct {
	array     unsafe.Pointer
	elemSize  uintptr
	size, i   int
	converter func(From) To
}

var _ c.Iterator[any] = (*ConvertIter[any, any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *ConvertIter[From, To]) For(walker func(element To) error) error {
	return loop.For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *ConvertIter[From, To]) ForEach(walker func(element To)) {
	loop.ForEach(i.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *ConvertIter[From, To]) Next() (To, bool) {
	if i.i < i.size {
		v := *(*From)(notsafe.GetArrayElemRef(i.array, i.i, i.elemSize))
		i.i++
		return i.converter(v), true
	}
	var no To
	return no, false
}

// Cap returns the iterator capacity
func (i *ConvertIter[From, To]) Cap() int {
	return i.size
}

func (i *ConvertIter[From, To]) Start() (*ConvertIter[From, To], To, bool) {
	return startIt[To](i)
}

func nextFiltered[T any](array unsafe.Pointer, size int, elemSize uintptr, filter func(T) bool, index *int) (T, bool) {
	for i := *index; i < size; i++ {
		if v := *(*T)(notsafe.GetArrayElemRef(array, i, elemSize)); filter(v) {
			*index = i + 1
			return v, true
		}
	}
	var v T
	return v, false
}
