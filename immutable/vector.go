package immutable

import (
	"fmt"
	"sort"

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

func (s *Vector[t]) Begin() c.Iterator[t] {
	return s.Iter()
}

func (s *Vector[t]) Iter() *it.Iter[t] {
	return it.New(s.elements)
}

func (s *Vector[t]) Collect() []t {
	return slice.Copy(s.elements)
}

func (s *Vector[t]) Len() int {
	return it.GetLen(&s.elements)
}

func (s *Vector[t]) Get(index int) (t, bool) {
	return slice.Get(s.elements, index)
}

func (s *Vector[t]) Track(tracker func(int, t) error) error {
	return slice.Track(s.elements, tracker)
}

func (s *Vector[t]) TrackEach(tracker func(int, t)) {
	slice.TrackEach(s.elements, tracker)
}

func (s *Vector[t]) For(walker func(t) error) error {
	return slice.For(s.elements, walker)
}

func (s *Vector[t]) ForEach(walker func(t)) {
	slice.ForEach(s.elements, walker)
}

func (s *Vector[t]) Filter(filter c.Predicate[t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Filter(s.Iter(), filter))
}

func (s *Vector[t]) Map(by c.Converter[t, t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Map(s.Iter(), by))
}

func (s *Vector[t]) Reduce(by op.Binary[t]) t {
	return it.Reduce(s.Iter(), by)
}

func (s *Vector[t]) Sort(less func(e1, e2 t) bool) *Vector[t] {
	var (
		elements = s.elements
		dest     = make([]t, len(elements))
	)
	copy(dest, elements)
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (s *Vector[t]) String() string {
	return slice.ToString(s.elements)
}
