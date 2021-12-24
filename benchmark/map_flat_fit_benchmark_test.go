package benchmark

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
)

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
		s = iter.ToSlice(iter.Map(iter.Filter(iter.Wrap(items), even), conv.And(toString, addTail)))
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
		s = iter.ToSlice(iter.Map(slice.Filter(items, even), conv.And(toString, addTail)))
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
		s = iter.ToSlice(iter.MapFit(iter.Wrap(items), even, conv.And(toString, addTail)))
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
		s = iter.ToSlice(slice.MapFit(items, even, conv.And(toString, addTail)))
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
		oneDimension := iter.ToSlice(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds))
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
		oneDimension := iter.ToSlice(iter.Filter(iter.Flatt(slice.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds))
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
		_ = iter.Reduce(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
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
		_ = iter.Reduce(iter.Filter(iter.Flatt(slice.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		oddSum := 0
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
	b.StopTimer()
}

func Benchmark_MapFlattStructure_Iterable(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice(iter.Map(iter.NotNil(iter.Flatt(iter.NotNil(iter.Wrap(items)), getAttributes)), getName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_Iterable_2(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := iter.Wrap([]*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice(iter.Map(iter.NotNil(iter.Flatt(iter.NotNil(items), getAttributes)), getName))
		items.(*iter.Slice[*Item]).Reset()
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFit(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice(iter.MapFit(iter.FlattFit(iter.Wrap(items), check.NotNil[*Item], getAttributes), check.NotNil[*Attributes], getName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFit2(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := iter.Wrap([]*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice(iter.MapFit(iter.FlattFit(items, check.NotNil[*Item], getAttributes), check.NotNil[*Attributes], getName))
		items.(*iter.Slice[*Item]).Reset()
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice(iter.MapFit(slice.FlattFit(items, check.NotNil[*Item], getAttributes), check.NotNil[*Attributes], getName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_Slice_PlainOld(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {

		names := make([]string, 0)
		for _, i := range items {
			if check.NotNil(i) {
				for _, a := range getAttributes(i) {
					if check.NotNil(a) {
						names = append(names, getName(a))
					}
				}
			}
		}
	}
	b.StopTimer()
}
