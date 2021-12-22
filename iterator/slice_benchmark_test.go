package iterator

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/conv"
)

var amount = 100_000

func Benchmark_WrapSliceHasNextGet(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := Wrap(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_WrapSliceNextNext(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := Wrap(values)
		for it.HasNext() {
			result = append(result, it.Get())
		}
	}
	b.StopTimer()
}

func Benchmark_NewSliceNextNext(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		it := New(values...)
		for v, ok := it.Next(); ok; v, ok = it.Next() {
			result = append(result, v)
		}
	}
	b.StopTimer()
}

func Benchmark_WrapSliceForEach(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		ForEach(Wrap(values), func(v int) { result = append(result, v) })
	}
	b.StopTimer()
}

func Benchmark_NewSliceHasNextGet(b *testing.B) {
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
		items    = New(1, 2, 3, 4, 5)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := SliceOf(Map[int](items, conv.And(toString, addTail), even))
		_ = s
	}
	b.StopTimer()
}

func Benchmark_Iterable_Flatt(b *testing.B) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = Of(Of(Of(1, 2, 3), Of(4, 5, 6)), Of(Of(7), nil), nil)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := SliceOf(Filter(Flatt(Flatt(multiDimension, conv.To[Iterator[Iterator[int]]]), conv.To[Iterator[int]]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}


func Benchmark_Iterable_Flatt2(b *testing.B) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension  = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := SliceOf(Filter(Flatt(Flatt(Iter(multiDimension), Iter[[]int]), Iter[int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

// func Benchmark_MapFlattDeepStructure(b *testing.B) {
// 	type (
// 		Attributes struct{ name string }
// 		Item       struct{ attributes []*Attributes }
// 	)

// 	var (
// 		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

// 		getName       = func(a *Attributes) string { return a.name }
// 		getAttributes = func(item *Item) []*Attributes { return item.attributes }
// 	)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = Map(Flatt(items, getAttributes, check.NotNil[*Item]), getName, check.NotNil[*Attributes])
// 	}
// 	b.StopTimer()
// }

// func Benchmark_MapFlattDeepStructure2(b *testing.B) {
// 	type (
// 		Attributes struct{ name string }
// 		Item       struct{ attributes []*Attributes }
// 	)

// 	var (
// 		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

// 		getName = func(a *Attributes) string { return a.name }
// 	)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = Flatt(items, func(item *Item) []string { return Map(item.attributes, getName, check.NotNil[*Attributes]) }, check.NotNil[*Item])

// 	}
// 	b.StopTimer()
// }

// func Benchmark_MapFlattDeepStructurePlainOld(b *testing.B) {
// 	type (
// 		Attributes struct{ name string }
// 		Item       struct{ attributes []*Attributes }
// 	)

// 	var (
// 		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

// 		getName       = func(a *Attributes) string { return a.name }
// 		getAttributes = func(item *Item) []*Attributes { return item.attributes }
// 	)
// 	b.ResetTimer()
// 	for j := 0; j < b.N; j++ {
// 		names := make([]string, 0)
// 		for _, i := range items {
// 			if check.NotNil(i) {
// 				for _, a := range getAttributes(i) {
// 					if check.NotNil(a) {
// 						names = append(names, getName(a))
// 					}
// 				}
// 			}
// 		}
// 	}
// 	b.StopTimer()
// }
