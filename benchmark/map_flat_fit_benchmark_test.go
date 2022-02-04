package benchmark

import (
	"fmt"
	"testing"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/c/op"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/conv"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it"
	iterimpl "github.com/m4gshm/gollections/it/impl/it"
	sliceitimpl "github.com/m4gshm/gollections/it/impl/slice"
	sliceit "github.com/m4gshm/gollections/it/slice"
	"github.com/m4gshm/gollections/mutable"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/sum"
)

var (
	toString = func(i int) string { return fmt.Sprintf("%d", i) }
	addTail  = func(s string) string { return s + "_tail" }
	even     = func(v int) bool { return v%2 == 0 }
)

func Benchmark_Map_EmbeddedSlice(b *testing.B) {
	op := conv.And(toString, addTail)
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

func Benchmark_Map_Iterator(b *testing.B) {
	op := conv.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items := it.Of(values...)
		s = it.Slice(it.Map(items, op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Iterator_Impl(b *testing.B) {
	op := conv.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items := iterimpl.New(values)
		s = iterimpl.Slice[string](iterimpl.Map(items, op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator(b *testing.B) {
	concat := conv.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.Slice(op.Map(items, concat))
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator_No_Cache_Operation(b *testing.B) {
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.Slice(op.Map(items, conv.And(toString, addTail)))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator_Impl(b *testing.B) {
	op := conv.And(toString, addTail)
	items := vector.New(values)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterimpl.Slice[string](iterimpl.Map(items.Iter(), op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Vector_ForEach_(b *testing.B) {
	op := conv.And(toString, addTail)
	items := vector.Of(values...)
	c := len(values)
	var s *mutable.Vector[string]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = mvector.New[string](c)
		items.ForEach(func(element int) { _ = s.Add(op(element)) })
	}
	_ = s
	b.StopTimer()
}

func Benchmark_MapAndFilter_Iterable(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.Slice(it.Map(it.Filter(it.Wrap(items), even), conv.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_MapAndFilter_Iterable_Impl(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterimpl.Slice[string](iterimpl.Map(iterimpl.Filter(iterimpl.New(items), even), conv.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_MapAndFilter_Slice(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.Slice(it.Map(sliceit.Filter(items, even), conv.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_MapAndFilter_Slice_Impl(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterimpl.Slice[string](iterimpl.Map(sliceitimpl.Filter(items, even), conv.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_MapFit_Iterable(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.Slice(it.MapFit(it.Wrap(items), even, conv.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_MapFit_Slice(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.Slice(sliceit.MapFit(items, even, conv.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_MapAndFilter_Slice_PlainOld(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
	)

	items := []int{1, 2, 3, 4, 5}
	var s []string
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		s := make([]string, 0)
		for _, i := range items {
			if i%2 == 0 {
				s = append(s, conv.And(toString, addTail)(i))
			}
		}
	}
	_ = s
	// fmt.Println(b.Name(), s)
	b.StopTimer()
}

func Benchmark_Flatt_Iterable(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := it.Slice(it.Filter(it.Flatt(it.Flatt(it.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := it.Slice(it.Filter(it.Flatt(sliceit.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		oneDimension := make([]int, 0)
		for _, i := range multiDimension {
			for _, ii := range i {
				for _, iii := range ii {
					if iii%2 != 0 {
						oneDimension = append(oneDimension, iii)
					}
				}
			}
		}
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Iterable(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Reduce(it.Filter(it.Flatt(it.Flatt(it.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds), sum.Of[int])
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Iterable_Impl(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iterimpl.Reduce(iterimpl.Filter(iterimpl.Flatt(iterimpl.Flatt(iterimpl.New(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds), sum.Of[int])
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Reduce(it.Filter(it.Flatt(sliceit.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), sum.Of[int])
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice_Impl(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iterimpl.Reduce(iterimpl.Filter(iterimpl.Flatt(sliceitimpl.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), sum.Of[int])
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	oddSum := 0
	for j := 0; j < b.N; j++ {
		oddSum = 0
		for _, i := range multiDimension {
			for _, ii := range i {
				for _, iii := range ii {
					if iii%2 != 0 {
						oddSum += iii
					}
				}
			}
		}
	}
	_ = oddSum
	b.StopTimer()
}

type (
	Attributes  struct{ name string }
	Participant struct{ attributes []*Attributes }
)

func (a *Attributes) GetName() string {
	if a == nil {
		return ""
	}
	return a.name
}

func (p *Participant) GetAttributes() []*Attributes {
	if p == nil {
		return nil
	}
	return p.attributes
}

func Benchmark_MapFlattStructure_IterableNotNil(b *testing.B) {

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Slice(it.Map(it.NotNil(it.Flatt(it.NotNil(it.Wrap(items)), (*Participant).GetAttributes)), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableWithoutNotNilFiltering(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Slice(it.Map(it.Flatt(it.Wrap(items), (*Participant).GetAttributes), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_Iterable_2(b *testing.B) {
	items := it.Wrap([]*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Slice(it.Map(it.NotNil(it.Flatt(it.NotNil(items), (*Participant).GetAttributes)), (*Attributes).GetName))
		items.(c.Resetable).Reset()
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Slice(it.MapFit(it.FlattFit(it.Wrap(items), check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFitReset_Impl(b *testing.B) {
	items := iterimpl.NewReseteable([]*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iterimpl.Slice[string](iterimpl.MapFit(iterimpl.FlattFit(items, check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
		items.Reset()
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Slice(it.MapFit(sliceit.FlattFit(items, check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit_Impl(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iterimpl.Slice[string](iterimpl.MapFit(sliceitimpl.FlattFit(items, check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceWithoutNilCheck(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.Slice(it.Map(sliceit.Flatt(items, (*Participant).GetAttributes), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_Slice_PlainOld(b *testing.B) {

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {

		names := make([]string, 0)
		for _, i := range items {
			if check.NotNil(i) {
				for _, a := range i.GetAttributes() {
					if check.NotNil(a) {
						names = append(names, a.GetName())
					}
				}
			}
		}
	}
	b.StopTimer()
}
