package iterator

type Iterator[T any] interface {
	Next() bool
	Get() T
}

func Wrap[T any](values []T) Iterator[T] {
	return &Slice[T]{values: values}
}

func New[T any](values ...T) Iterator[T] {
	return &Slice[T]{values: values}
}

func WrapMap[K comparable, V any](values map[K]V) Iterator[*KV[K, V]] {
	return newMap(values)
}
