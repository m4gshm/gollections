package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// NewVector instantiates Vector and copies elements to it.
func NewVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// WrapVector instantiates Vector using a slise as internal storage.
func WrapVector[T any](elements []T) *Vector[T] {
	return &Vector[T]{elements: elements, esize: notsafe.GetTypeSize[T]()}
}

// Vector is the Collection implementation that provides elements order and index access.
type Vector[T any] struct {
	elements []T
	esize    uintptr
}

var (
	_ c.Vector[any] = (*Vector[any])(nil)
	_ fmt.Stringer  = (*Vector[any])(nil)
)

func (v *Vector[T]) Begin() c.Iterator[T] {
	h := v.Head()
	return &h
}

func (v *Vector[T]) Head() iter.ArrayIter[T] {
	var (
		elements []T
		esize    uintptr
	)
	if v != nil {
		elements = v.elements
		esize = v.esize

	}
	return iter.NewHeadS(elements, esize)
}

func (v *Vector[T]) Tail() iter.ArrayIter[T] {
	var (
		elements []T
		esize    uintptr
	)
	if v != nil {
		elements = v.elements
		esize = v.esize

	}
	return iter.NewTailS(elements, esize)
}

func (v *Vector[T]) First() (iter.ArrayIter[T], T, bool) {
	var (
		elements []T
		esize    uintptr
	)
	if v != nil {
		elements = v.elements
		esize = v.esize

	}
	var (
		iterator  = iter.NewHeadS(elements, esize)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (v *Vector[T]) Last() (iter.ArrayIter[T], T, bool) {
	var (
		iterator  = iter.NewTailS(v.elements, v.esize)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

func (v *Vector[T]) Slice() (out []T) {
	if v == nil {
		return
	}
	return slice.Clone(v.elements)
}

func (v *Vector[T]) Len() int {
	if v == nil {
		return 0
	}
	return notsafe.GetLen(v.elements)
}

func (v *Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

func (v *Vector[T]) Get(index int) (out T, ok bool) {
	if v == nil {
		return
	}
	return slice.Get(v.elements, index)
}

func (v *Vector[T]) Track(tracker func(int, T) error) error {
	if v == nil {
		return nil
	}
	return slice.Track(v.elements, tracker)
}

func (v *Vector[T]) TrackEach(tracker func(int, T)) {
	if v != nil {
		slice.TrackEach(v.elements, tracker)
	}
}

func (v *Vector[T]) For(walker func(T) error) error {
	if v == nil {
		return nil
	}
	return slice.For(v.elements, walker)
}

func (v *Vector[T]) ForEach(walker func(T)) {
	if v != nil {
		slice.ForEach(v.elements, walker)
	}
}

func (v *Vector[T]) Filter(filter func(T) bool) c.Pipe[T] {
	h := v.Head()
	return iter.NewPipe[T](iter.Filter(h, h.Next, filter))
}

func (v *Vector[T]) Convert(by func(T) T) c.Pipe[T] {
	h := v.Head()
	return iter.NewPipe[T](iter.Convert(h, h.Next, by))
}

func (v *Vector[T]) Reduce(by func(T, T) T) T {
	h := v.Head()
	return loop.Reduce(h.Next, by)
}

func (v *Vector[T]) Sort(less slice.Less[T]) *Vector[T] {
	return v.sortBy(sort.Slice, less)
}

func (v *Vector[T]) StableSort(less slice.Less[T]) *Vector[T] {
	return v.sortBy(sort.SliceStable, less)
}

func (v *Vector[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *Vector[T] {
	var elements []T
	if v != nil {
		elements = slice.Clone(v.elements)
	}
	slice.Sort(elements, sorter, less)
	return WrapVector(elements)
}

func (v *Vector[T]) String() string {
	if v == nil {
		return ""
	}
	return slice.ToString(v.elements)
}
