package collection

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
)

// Seq is an alias of an iterator-function that allows to iterate over elements of a sequence, such as slice.
type Seq[T any] func(yield func(T) bool)

// SeqE is a specific iterator form that allows to retrieve a value with an error as second parameter of the iterator.
// It is used as a result of applying functions like seq.Conv, which may throw an error during iteration.
// At each iteration step, it is necessary to check for the occurrence of an error.
//
//	for e, err := range seqence {
//	    if err != nil {
//	        break
//	    }
//	    ...
//	}
type SeqE[T any] func(yield func(T, error) bool)

// Seq2 is an alias of an iterator-function that allows to iterate over key/value pairs of a sequence, such as slice or map.
// It is used to iterate over slice index/value pairs or map key/value pairs.
type Seq2[K, V any] func(yield func(K, V) bool)

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

func (s Seq2[K, V]) FilterKey(predicate predicate.Predicate[K]) Seq2[K, V] {
	return seq2.Filter(s, func(k K, v V) bool {
		return predicate(k)
	})
}

func (s Seq2[K, V]) ConvertValue(convertr func(v V) V) Seq2[K, V] {
	return seq2.Convert(s, func(k K, v V) (K, V) {
		return k, convertr(v)
	})
}
