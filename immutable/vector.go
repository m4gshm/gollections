package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// NewVector instantiates Vector and copies elements to it.
func NewVector[T any](elements []T) Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// WrapVector instantiates Vector using a slise as internal storage.
func WrapVector[T any](elements []T) Vector[T] {
	return Vector[T]{elements: elements, esize: notsafe.GetTypeSize[T]()}
}

// Vector is the Collection implementation that provides elements order and index access.
type Vector[T any] struct {
	elements []T
	esize    uintptr
}

var (
	_ c.Vector[any] = (*Vector[any])(nil)
	_ c.Vector[any] = Vector[any]{}
	_ fmt.Stringer  = (*Vector[any])(nil)
	_ fmt.Stringer  = Vector[any]{}
)

func (v Vector[T]) Begin() c.Iterator[T] {
	h := v.Head()
	return &h
}

func (v Vector[T]) Head() it.ArrayIter[T] {
	return it.NewHeadS(v.elements, v.esize)
}

func (v Vector[T]) Tail() it.ArrayIter[T] {
	return it.NewTailS(v.elements, v.esize)
}

func (v Vector[T]) First() (it.ArrayIter[T], T, bool) {
	var (
		iter      = it.NewHeadS(v.elements, v.esize)
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (v Vector[T]) Last() (it.ArrayIter[T], T, bool) {
	var (
		iter      = it.NewTailS(v.elements, v.esize)
		first, ok = iter.Prev()
	)
	return iter, first, ok
}

func (v Vector[T]) Collect() []T {
	return slice.Clone(v.elements)
}

func (v Vector[T]) Len() int {
	return notsafe.GetLen(v.elements)
}

func (v Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

func (v Vector[T]) Get(index int) (T, bool) {
	return slice.Get(v.elements, index)
}

func (v Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(v.elements, tracker)
}

func (v Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(v.elements, tracker)
}

func (v Vector[T]) For(walker func(T) error) error {
	return slice.For(v.elements, walker)
}

func (v Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(v.elements, walker)
}

func (v Vector[T]) Filter(filter func(T) bool) c.Pipe[T, []T] {
	h := v.Head()
	return it.NewPipe[T](it.Filter(h, h.Next, filter))
}

func (v Vector[T]) Convert(by func(T) T) c.Pipe[T, []T] {
	h := v.Head()
	return it.NewPipe[T](it.Convert(h, h.Next, by))
}

func (v Vector[T]) Reduce(by func(T, T) T) T {
	h := v.Head()
	return loop.Reduce(h.Next, by)
}

func (v Vector[T]) Sort(less slice.Less[T]) Vector[T] {
	return v.sortBy(sort.Slice, less)
}

func (v Vector[T]) StableSort(less slice.Less[T]) Vector[T] {
	return v.sortBy(sort.SliceStable, less)
}

func (v Vector[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) Vector[T] {
	c := slice.Clone(v.elements)
	slice.Sort(c, sorter, less)
	return WrapVector(c)
}

func (v Vector[T]) String() string {
	return slice.ToString(v.elements)
}
