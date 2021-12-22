package slice

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
)

var amount = 100_000

func Benchmark_SliceForEach(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, 100000)
		add := func(v int) { result = append(result, i) }
		ForEach(values, add)
	}
	b.StopTimer()
}

func Benchmark_EmbeddedSliceForAppendByFunc(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, amount)
		add := func(v int) { result = append(result, i) }
		for _, v := range values {
			add(v)
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedSliceForByIndexAppendByFunc(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, amount)
		add := func(v int) { result = append(result, i) }
		for ii := 0; ii < len(values); ii++ {
			v := values[ii]
			add(v)
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedSliceFor(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, amount)
		for _, v := range values {
			result = append(result, v)
		}
	}
	b.StopTimer()
}

func Benchmark_EmbeddedSliceForByIndex(b *testing.B) {
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, amount)
		for ii := 0; ii < len(values); ii++ {
			v := values[ii]
			result = append(result, v)
		}
	}
	b.StopTimer()
}

func Benchmark_MapAndFilter(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
		items    = Of(1, 2, 3, 4, 5)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(items, conv.And(toString, addTail), even)
	}
	b.StopTimer()
}

func Benchmark_MapAndFilterPlainOld(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
		items    = Of(1, 2, 3, 4, 5)
	)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		converted := make([]string, 0)
		for _, i := range items {
			if even(i) {
				converted = append(converted, conv.And(toString, addTail)(i))
			}
		}
		_ = converted
	}
	b.StopTimer()
}

func Benchmark_FlattSlices(b *testing.B) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(Flatt(Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds)
	}
	b.StopTimer()
}

func Benchmark_FlattSlices2(b *testing.B) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Flatt(multiDimension, func(i2 [][]int) []int { return Flatt(i2, func(i1 []int) []int { return Filter(i1, odds) }) })
	}
	b.StopTimer()
}

func Benchmark_FlattSlicesPlainOld(b *testing.B) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		oneDimension := make([]int, 0)
		for _, i := range multiDimension {
			for _, ii := range i {
				for _, iii := range ii {
					if odds(iii) {
						oneDimension = append(oneDimension, iii)
					}
				}
			}
		}
	}
	b.StopTimer()
}

func Benchmark_MapFlattDeepStructure(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(Flatt(items, getAttributes, check.NotNil[*Item]), getName, check.NotNil[*Attributes])
	}
	b.StopTimer()
}

func Benchmark_MapFlattDeepStructure2(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

		getName = func(a *Attributes) string { return a.name }
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Flatt(items, func(item *Item) []string { return Map(item.attributes, getName, check.NotNil[*Attributes]) }, check.NotNil[*Item])

	}
	b.StopTimer()
}

func Benchmark_MapFlattDeepStructurePlainOld(b *testing.B) {
	type (
		Attributes struct{ name string }
		Item       struct{ attributes []*Attributes }
	)

	var (
		items = []*Item{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

		getName       = func(a *Attributes) string { return a.name }
		getAttributes = func(item *Item) []*Attributes { return item.attributes }
	)
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
