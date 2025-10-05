package seq

// Slice collects the elements of the 'seq' sequence into a new slice.
func (s Seq[T]) Slice() []T {
	return Slice(s)
}

// Append collects the elements of the 'seq' sequence into the specified 'out' slice.
func (s Seq[T]) Append(out []T) []T {
	return Append(s, out)
}

// Reduce reduces the elements of the seq into one using the 'merge' function.
func (s Seq[T]) Reduce(merge func(a T, b T) T) T {
	return Reduce(s, merge)
}

// ReduceOK reduces the elements of the seq into one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func (s Seq[T]) ReduceOK(merge func(T, T) T) (result T, ok bool) {
	return ReduceOK(s, merge)
}

// Reducee reduces the elements of the seq into one using the 'merge' function.
func (s Seq[T]) Reducee(merge func(T, T) (T, error)) (T, error) {
	return Reducee(s, merge)
}

// ReduceeOK reduces the elements of the seq into one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func (s Seq[T]) ReduceeOK(merge func(T, T) (T, error)) (result T, ok bool, err error) {
	return ReduceeOK(s, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func (s Seq[T]) Accum(first T, merge func(T, T) T) T {
	return Accum(first, s, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func (s Seq[T]) Accumm(first T, merge func(T, T) (T, error)) (T, error) {
	return Accumm(first, s, merge)
}

// Head returns the first element.
func (s Seq[T]) Head() (v T, ok bool) {
	return Head(s)
}

// First returns the first element that satisfies the condition.
func (s Seq[T]) First(predicate func(T) bool) (v T, ok bool) {
	return First(s, predicate)
}

// Firstt returns the first element that satisfies the condition.
func (s Seq[T]) Firstt(predicate func(T) (bool, error)) (v T, ok bool, err error) {
	return Firstt(s, predicate)
}

// Top returns a sequence of top n elements.
func (s Seq[T]) Top(n int) Seq[T] {
	return Top(n, s)
}

// Skip returns the seq without first n elements.
func (s Seq[T]) Skip(n int) Seq[T] {
	return Skip(n, s)
}

// While cuts tail elements of the seq that don't match the filter.
func (s Seq[T]) While(filter func(T) bool) Seq[T] {
	return While(s, filter)
}

// SkipWhile returns a sequence without first elements of the seq that dont'math the filter.
func (s Seq[T]) SkipWhile(filter func(T) bool) Seq[T] {
	return SkipWhile(s, filter)
}

// HasAny checks whether the seq contains an element that satisfies the condition.
func (s Seq[T]) HasAny(predicate func(T) bool) bool {
	return HasAny(s, predicate)
}

// Union combines several sequences into one.
func (s Seq[T]) Union(seqences ...seq[T]) Seq[T] {
	return Union(append(append(make([]seq[T], len(seqences)+1), s), seqences...)...)
}

// Filter creates an iterator that iterates only those elements for which the 'filter' function returns true.
func (s Seq[T]) Filter(filter func(s T) bool) Seq[T] {
	return Filter(s, filter)
}

// Filt creates an erroreable iterator that iterates only those elements for which the 'filter' function returns true.
func (s Seq[T]) Filt(filter func(s T) (bool, error)) SeqE[T] {
	return Filt(s, filter)
}

// Convert creates an iterator that applies the 'converter' function to each iterable element.
func (s Seq[T]) Convert(converter func(t T) T) Seq[T] {
	return Convert(s, converter)
}

// Conv creates an errorable seq that applies the 'converter' function to the iterable elements.
func (s Seq[T]) Conv(converter func(T) (T, error)) SeqE[T] {
	return Conv(s, converter)
}

// ForEach applies the 'consumer' function to the seq elements.
func (s Seq[T]) ForEach(f func(T)) {
	ForEach(s, f)
}
