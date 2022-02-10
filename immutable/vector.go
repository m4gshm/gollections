package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewVector[t any](elements []t) *Vector[t] {
	return WrapVector(slice.Copy(elements))
}

func WrapVector[t any](elements []t) *Vector[t] {
	return &Vector[t]{elements: elements}
}

//Vector stores ordered elements, provides index access.
type Vector[t any] struct {
	elements []t
}

var (
	_ c.Vector[any] = (*Vector[any])(nil)
	_ fmt.Stringer  = (*Vector[any])(nil)
)

func (v *Vector[t]) Begin() c.Iterator[t] {
	return v.Iter()
}

func (v *Vector[t]) Iter() *it.Iter[t] {
	return it.New(v.elements)
}

func (v *Vector[t]) Collect() []t {
	return slice.Copy(v.elements)
}

func (v *Vector[t]) Len() int {
	return it.GetLen(&v.elements)
}

func (v *Vector[t]) Get(index int) (t, bool) {
	return slice.Get(v.elements, index)
}

func (v *Vector[t]) Track(tracker func(int, t) error) error {
	return slice.Track(v.elements, tracker)
}

func (v *Vector[t]) TrackEach(tracker func(int, t)) {
	slice.TrackEach(v.elements, tracker)
}

func (v *Vector[t]) For(walker func(t) error) error {
	return slice.For(v.elements, walker)
}

func (v *Vector[t]) ForEach(walker func(t)) {
	slice.ForEach(v.elements, walker)
}

func (v *Vector[t]) Filter(filter c.Predicate[t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Filter(v.Iter(), filter))
}

func (v *Vector[t]) Map(by c.Converter[t, t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Map(v.Iter(), by))
}

func (v *Vector[t]) Reduce(by op.Binary[t]) t {
	return it.Reduce(v.Iter(), by)
}

func (v *Vector[t]) Sort(less func(e1, e2 t) bool) *Vector[t] {
	return WrapVector(slice.SortCopy(v.elements, less))
}

func (v *Vector[t]) String() string {
	return slice.ToString(v.elements)
}
