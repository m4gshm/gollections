package immutable

import "fmt"

func NewSet[T comparable](values ...T) Set[T] {
	result := Set[T]{}
	for _, a := range values {
		result[a] = struct{}{}
	}
	return result
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for k := range s {
		values = append(values, k)

	}
	return values
}

func (s Set[T]) String() string {
	str := ""
	for k := range s {
		if len(str) > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%+v", k)
	}
	return "[" + str + "]"
}

func (s Set[T]) Len() int {
	return len(s)
}

var _ fmt.Stringer = (Set[interface{}])(nil)
