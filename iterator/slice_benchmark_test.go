package iterator

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

var amount = 100_000

func Benchmark_IterateSlice_NextGet(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := Iter(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_IterateSlice_ForEach(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		ForEach(Iter(values), func(v int) { result = append(result, v) })
	}
	b.StopTimer()
}

func Benchmark_NewSlice_ForEach(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		ForEach(New(values...), func(v int) { result = append(result, v) })
	}
	b.StopTimer()
}

func Benchmark_NewSlice_HasNextGet(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := New(values...)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_Iterable_MapAndFilter(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := New(1, 2, 3, 4, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := ToSlice(Map[int](items, conv.And(toString, addTail), even))
		_ = s
	}
	b.StopTimer()
}

func Benchmark_Iterable_Flatt(b *testing.B) {
	var (
		odds = func(v int) bool { return v%2 != 0 }
	)
	multiDimension := Iter([][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := ToSlice(Filter(Flatt(Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Iterable_MapFlattDeepStructure(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
	items := Iter([]*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToSlice(Map(Flatt(items, getAttributes, check.NotNil[*Item]), getName, check.NotNil[*Attributes]))
	}
	b.StopTimer()
}
