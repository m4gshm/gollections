package seq

import "github.com/m4gshm/gollections/internal/seqe"

// Slice collects the elements of the 'seq' sequence into a new slice.
func (s SeqE[T]) Slice() ([]T, error) {
	return seqe.Slice(s)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func (s SeqE[T]) Append(out []T) ([]T, error) {
	return seqe.Append(s, out)
}

// Reduce reduces the elements of the seq into one using the 'merge' function.
func (s SeqE[T]) Reduce(merge func(a T, b T) T) (T, error) {
	return seqe.Reduce(s, merge)
}

// ReduceOK reduces the elements of the seq into one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func (s SeqE[T]) ReduceOK(merge func(T, T) T) (result T, ok bool, err error) {
	return seqe.ReduceOK(s, merge)
}

// Reducee reduces the elements of the seq into one using the 'merge' function.
func (s SeqE[T]) Reducee(merge func(T, T) (T, error)) (T, error) {
	return seqe.Reducee(s, merge)
}

// ReduceeOK reduces the elements of the seq into one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func (s SeqE[T]) ReduceeOK(merge func(T, T) (T, error)) (result T, ok bool, err error) {
	return seqe.ReduceeOK(s, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func (s SeqE[T]) Accum(first T, merge func(T, T) T) (T, error) {
	return seqe.Accum(first, s, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func (s SeqE[T]) Accumm(first T, merge func(T, T) (T, error)) (T, error) {
	return seqe.Accumm(first, s, merge)
}

// Head returns the first element.
func (s SeqE[T]) Head() (T, bool, error) {
	return seqe.Head(s)
}

// First returns the first element that satisfies the condition.
func (s SeqE[T]) First(predicate func(T) bool) (T, bool, error) {
	return seqe.First(s, predicate)
}

// Firstt returns the first element that satisfies the condition.
func (s SeqE[T]) Firstt(predicate func(T) (bool, error)) (T, bool, error) {
	return seqe.Firstt(s, predicate)
}

// Top returns a sequence of top n elements.
func (s SeqE[T]) Top(n int) SeqE[T] {
	return seqe.Top(n, s)
}

// Skip returns the seq without first n elements.
func (s SeqE[T]) Skip(n int) SeqE[T] {
	return seqe.Skip(n, s)
}

// While cuts tail elements of the seq that don't match the filter.
func (s SeqE[T]) While(filter func(T) bool) SeqE[T] {
	return seqe.While(s, filter)
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func (s SeqE[T]) SkipWhile(filter func(T) bool) SeqE[T] {
	return seqe.SkipWhile(s, filter)
}

// HasAny checks whether the seq contains an element that satisfies the condition.
func (s SeqE[T]) HasAny(predicate func(T) bool) (bool, error) {
	return seqe.HasAny(s, predicate)
}

// Union combines several sequences into one.
func (s SeqE[T]) Union(seqences ...seqE[T]) SeqE[T] {
	return seqe.Union(append(append(make([]seqE[T], len(seqences)+1), s), seqences...)...)
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func (s SeqE[T]) Filter(filter func(s T) bool) SeqE[T] {
	return seqe.Filter(s, filter)
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func (s SeqE[T]) Filt(filter func(s T) (bool, error)) SeqE[T] {
	return seqe.Filt(s, filter)
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func (s SeqE[T]) Convert(converter func(t T) T) SeqE[T] {
	return seqe.Convert(s, converter)
}

// Conv creates an errorable seq that applies the 'converter' function to the collection elements.
func (s SeqE[T]) Conv(converter func(T) (T, error)) SeqE[T] {
	return seqe.Conv(s, converter)
}

// ForEach applies the 'consumer' function to the seq elements.
func (s SeqE[T]) ForEach(f func(T)) error {
	return seqe.ForEach(s, f)
}
