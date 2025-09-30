package seq

func (s Seq[T]) Slice() []T {
	return Slice(s)
}

func (s Seq[T]) Append(out []T) []T {
	return Append(s, out)
}

func (s Seq[T]) Reduce(merge func(a T, b T) T) T {
	return Reduce(s, merge)
}

// ReduceOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func (s Seq[T]) ReduceOK(merge func(T, T) T) (result T, ok bool) {
	return ReduceOK(s, merge)
}

func (s Seq[T]) Reducee(merge func(T, T) (T, error)) (T, error) {
	return Reducee(s, merge)
}

// ReduceeOK reduces the elements of the 'seq' sequence an one using the 'merge' function.
// Returns ok==false if the seq returns ok=false at the first call (no more elements).
func (s Seq[T]) ReduceeOK(merge func(T, T) (T, error)) (result T, ok bool, err error) {
	return ReduceeOK(s, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element of the 'seq' sequence.
func (s Seq[T]) Accum(first T, merge func(T, T) T) T {
	return Accum(first, s, merge)
}

func (s Seq[T]) Head() (v T, ok bool) {
	return Head(s)
}

func (s Seq[T]) First(predicate func(T) bool) (v T, ok bool) {
	return First(s, predicate)
}

func (s Seq[T]) Firstt(predicate func(T) (bool, error)) (v T, ok bool, err error) {
	return Firstt(s, predicate)
}

func (s Seq[T]) Top(n int) Seq[T] {
	return Top(n, s)
}

func (s Seq[T]) Skip(n int) Seq[T] {
	return Skip(n, s)
}

func (s Seq[T]) HasAny(predicate func(T) bool) bool {
	return HasAny(s, predicate)
}

func (s Seq[T]) Union(seqences ...seq[T]) Seq[T] {
	return Union(append(append(make([]seq[T], len(seqences)+1), s), seqences...)...)
}

func (s Seq[T]) Filter(filter func(s T) bool) Seq[T] {
	return Filter(s, filter)
}

func (s Seq[T]) Filt(filter func(s T) (bool, error)) SeqE[T] {
	return Filt(s, filter)
}

func (s Seq[T]) Convert(converter func(t T) T) Seq[T] {
	return Convert(s, converter)
}

func (s Seq[T]) Conv(converter func(T) (T, error)) SeqE[T] {
	return Conv(s, converter)
}

func (s Seq[T]) ForEach(f func(T)) {
	ForEach(s, f)
}
