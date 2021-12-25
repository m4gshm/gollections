package immutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

func NewSet[T comparable](values ...T) *OrderSet[T] {
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
	return &OrderSet[T]{values: order}
}

type OrderSet[T comparable, S int] struct {
	values []T
}

var _ typ.Walk[interface{}] = (*OrderSet[interface{}])(nil)
var _ typ.Container[interface{}, int] = (*OrderSet[interface{}, int])(nil)
var _ fmt.Stringer = (*OrderSet[interface{}])(nil)

func (s *OrderSet[T,S]) Begin() typ.Iterator[T] {
	return iter.New(s.values)
}

func (s *OrderSet[T,S]) Values() []T {
	out := make([]T, len(s.values))
	copy(out, s.values)
	return out
}

func (s *OrderSet[T, S]) Len() int {
	return len(s.values)
}

func (s *OrderSet[T, S]) String() string {
	str := ""
	for _, v := range s.values {
		if len(str) > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%+v", v)
	}
	return "[" + str + "]"
}

