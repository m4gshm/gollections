package convert

import (
	"fmt"
	"testing"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/vector"
	"github.com/m4gshm/gollections/collection/mutable"
	mvector "github.com/m4gshm/gollections/collection/mutable/vector"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	toString = func(i int) string { return fmt.Sprintf("%d", i) }
	addTail  = func(s string) string { return s + "_tail" }
	even     = func(v int) bool { return v%2 == 0 }
	max      = 100000
	values   = range_.Closed(1, max)
)

func Benchmark_Convert_Slice(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = slice.Convert(values, op)
	}
	_ = len(s)
	b.StopTimer()
}

func Benchmark_Convert_Slice_EveryElement(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = make([]string, len(values))
		for i, v := range values {
			s[i] = op(v)
		}
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_Seq(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := seq.Of(values...)
		s = seq.SliceCap(seq.Convert(it, op), len(values))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_Loop(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := slice.NewHead(values)
		s = loop.SliceCap(loop.Convert(it.Next, op), len(values))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_ImmutableVector_Iterable(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = collection.Convert(items, concat).Slice()
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Convert_ImmutableVector_Iterable_Append(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = collection.Convert(items, concat).Append(make([]string, 0, len(values)))
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Convert_ImmutableVector(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = vector.Convert(items, concat).Slice()
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_ImmutableVector_Append(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = vector.Convert(items, concat).Append(make([]string, 0, len(values)))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_ImmutableVector_Head_Loop(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = loop.SliceCap(loop.Convert(ptr.Of(items.Head()).Next, concat), len(values))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_ImmutableVector_ForEach_To_MutableVector(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	c := len(values)
	var s *mutable.Vector[string]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = mvector.NewCap[string](c)
		items.ForEach(func(element int) { s.Add(concat(element)) })
	}
	_ = s
	b.StopTimer()
}
