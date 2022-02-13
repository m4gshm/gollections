package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Copy(elements))
}

func WrapVector[T any](elements []T) *Vector[T] {
	return &Vector[T]{elements: elements, esize: it.GetTypeSize[T]()}
}

//Vector stores ordered elements, provides index access.
type Vector[T any] struct {
	elements []T
	esize    uintptr
}

var (
	_ c.Vector[any] = (*Vector[any])(nil)
	_ fmt.Stringer  = (*Vector[any])(nil)
)

func (v *Vector[T]) Begin() c.Iterator[T] {
	return v.Head()
}

func (v *Vector[T]) Head() *it.Iter[T] {
	return it.NewHeadS(v.elements, v.esize)
}

func (v *Vector[T]) Tail() *it.Iter[T] {
	return it.NewTailS(v.elements, v.esize)
}

func (v *Vector[T]) Collect() []T {
	return slice.Copy(v.elements)
}

func (v *Vector[T]) Len() int {
	return it.GetLen(v.elements)
}

func (v *Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

func (v *Vector[T]) Get(index int) (T, bool) {
	return slice.Get(v.elements, index)
}

func (v *Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(v.elements, tracker)
}

func (v *Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(v.elements, tracker)
}

func (v *Vector[T]) For(walker func(T) error) error {
	return slice.For(v.elements, walker)
}

func (v *Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(v.elements, walker)
}

func (v *Vector[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Filter(v.Head(), filter))
}

func (v *Vector[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Map(v.Head(), by))
}

func (v *Vector[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(v.Head(), by)
}

func (v *Vector[T]) Sort(less func(e1, e2 T) bool) *Vector[T] {
	return WrapVector(slice.SortCopy(v.elements, less))
}

func (v *Vector[T]) String() string {
	return slice.ToString(v.elements)
}
