package slice

func Wrap[T any](elements ...T) Slice[T] { return elements }

type Slice[T any] []T

func (s Slice[T]) Unwrap() []T {
	return s
}

func (s Slice[T]) Filter(filter func(T) bool) Slice[T] {
	return Filter(s, filter)
}

func (s Slice[T]) AppendFilter(filter func(T) bool) Slice[T] {
	return AppendFilter(s, s, filter)
}

func (s Slice[T]) Filt(filter func(T) (bool, error)) (Slice[T], error) {
	return Filt(s, filter)
}

func (s Slice[T]) AppendFilt(filter func(T) (bool, error)) (Slice[T], error) {
	return AppendFilt(s, s, filter)
}

func (s Slice[T]) Clone() Slice[T] {
	return Clone(s)
}

func (s Slice[T]) DeepClone(copier func(T) T) Slice[T] {
	return DeepClone(s, copier)
}

func (s Slice[T]) Reduce(merge func(T, T) T) (out T) {
	return Reduce(s, merge)
}

func (s Slice[T]) First(by func(T) bool) (no T, ok bool) {
	return First(s, by)
}

func (s Slice[T]) Firstt(by func(T) (bool, error)) (no T, ok bool, err error) {
	return Firstt(s, by)
}

func (s Slice[T]) Last(by func(T) bool) (no T, ok bool) {
	return Last(s, by)
}

func (s Slice[T]) Lastt(by func(T) (bool, error)) (no T, ok bool, err error) {
	return Lastt(s, by)
}

func (s Slice[T]) HasAny( predicate func(T) bool) bool {
	return HasAny(s, predicate)
}

func (s Slice[T]) Empty() bool {
	return Empty(s)
}