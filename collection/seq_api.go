package collection

import "github.com/m4gshm/gollections/seq"

func (s Seq[T]) Slice() []T {
	return seq.Slice(s)
}

func (s Seq[T]) Reduce(merge func(a T, b T) T) T {
	return seq.Reduce(s, merge)
}

func (s Seq[T]) Filter(filter func(s T) bool) Seq[T] {
	return seq.Filter(s, filter)
}

func (s Seq[T]) HasAny(predicate func(T) bool) bool {
	return seq.HasAny(s, predicate)
}

func (s Seq[T]) Convert(converter func(t T) T) Seq[T] {
	return seq.Convert(s, converter)
}

func (s Seq[T]) Append(out []T) []T {
	return seq.Append(s, out)
}

func (s Seq[T]) ForEach(f func(T)) {
	seq.ForEach(s, f)
}
