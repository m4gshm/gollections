package immutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

func NewSet[T comparable](values ...T) *OrderedSet[T] {
	var (
		uniques = make(map[T]struct{}, 0)
		order   = make([]T, 0, 0)
	)
	for _, v := range values {
		if _, ok := uniques[v]; !ok {
			order = append(order, v)
			uniques[v] = struct{}{}
		}
	}
	return &OrderedSet[T]{values: order}
}

type OrderedSet[T comparable, S int] struct {
	values []T
}

var _ typ.Walk[interface{}] = (*OrderedSet[interface{}])(nil)
var _ typ.Container[interface{}, int] = (*OrderedSet[interface{}, int])(nil)
var _ fmt.Stringer = (*OrderedSet[interface{}])(nil)

func (s *OrderedSet[T, S]) Begin() typ.Iterator[T] {
	return iter.New(s.values)
}

func (s *OrderedSet[T, S]) Values() []T {
	out := make([]T, len(s.values))
	copy(out, s.values)
	return out
}

func (s *OrderedSet[T, S]) Len() int {
	return len(s.values)
}

func (s *OrderedSet[T, S]) String() string {
	str := ""
	for _, v := range s.values {
		if len(str) > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%+v", v)
	}
	return "[" + str + "]"
}
