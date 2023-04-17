package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys is non-copy constructor
func WrapKeys[K comparable, V any](uniques map[K]V) MapKeys[K, V] {
	return MapKeys[K, V]{uniques}
}

// MapKeys is the container reveal keys of a map and hides values.
type MapKeys[K comparable, V any] struct {
	uniques map[K]V
}

var (
	_ c.Collection[int, []int, c.Iterator[int]] = (*MapKeys[int, any])(nil)
	_ c.Collection[int, []int, c.Iterator[int]] = MapKeys[int, any]{}
	_ fmt.Stringer                              = (*MapValues[int, any])(nil)
	_ fmt.Stringer                              = MapValues[int, any]{}
)

func (s MapKeys[K, V]) Begin() c.Iterator[K] {
	h := s.Head()
	return &h
}

func (s MapKeys[K, V]) Head() it.Key[K, V] {
	return it.NewKey(s.uniques)
}

func (s MapKeys[K, V]) First() (it.Key[K, V], K, bool) {
	var (
		iter      = s.Head()
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (s MapKeys[K, V]) Len() int {
	return len(s.uniques)
}

func (s MapKeys[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s MapKeys[K, V]) Collect() []K {
	uniques := s.uniques
	elements := make([]K, 0, len(s.uniques))
	for key := range uniques {
		elements = append(elements, key)
	}
	return elements
}

func (s MapKeys[K, V]) For(walker func(K) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s MapKeys[K, V]) ForEach(walker func(K)) {
	map_.ForEachKey(s.uniques, walker)
}

func (s MapKeys[K, V]) Filter(filter func(K) bool) c.Pipe[K, []K] {
	h := s.Head()
	return it.NewPipe[K](it.Filter(h, h.Next, filter))
}

func (s MapKeys[K, V]) Convert(by func(K) K) c.Pipe[K, []K] {
	h := s.Head()
	return it.NewPipe[K](it.Convert(h, h.Next, by))
}

func (s MapKeys[K, V]) Reduce(by func(K, K) K) K {
	return loop.Reduce(s.Head().Next, by)
}

func (s MapKeys[K, V]) String() string {
	return slice.ToString(s.Collect())
}
