package benchmark

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/m4gshm/gollections/c/op"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/it"
	iterimpl "github.com/m4gshm/gollections/it/impl/it"
	sliceitimpl "github.com/m4gshm/gollections/it/impl/slice"
	sliceit "github.com/m4gshm/gollections/it/slice"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/mutable"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	sop "github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
)

var (
	toString = func(i int) string { return fmt.Sprintf("%d", i) }
	addTail  = func(s string) string { return s + "_tail" }
	even     = func(v int) bool { return v%2 == 0 }
)

func Benchmark_Map_EmbeddedSlice(b *testing.B) {
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

func Benchmark_First_Slice(b *testing.B) {
	op := func(i int) bool { return i > 10000 }
	var f int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, ok = slice.First(values, op)
	}
	_ = f
	_ = ok
	b.StopTimer()
}

func Benchmark_First_OfBy(b *testing.B) {
	op := func(i int) bool { return i > 10000 }
	var f int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, ok = first.Of(values...).By(op)
	}
	_ = f
	_ = ok
	b.StopTimer()
}

func Benchmark_First_Iterator(b *testing.B) {
	op := func(i int) bool { return i > 10000 }
	var f int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, ok = it.First(it.Of(values...), op)
	}
	_ = f
	_ = ok
	b.StopTimer()
}

func Benchmark_Last_Slice(b *testing.B) {
	op := func(i int) bool { return i < 10000 }
	var f int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, ok = slice.Last(values, op)
	}
	_ = f
	_ = ok
	b.StopTimer()
}

func Benchmark_Last_OfBy(b *testing.B) {
	op := func(i int) bool { return i < 10000 }
	var f int
	var ok bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, ok = last.Of(values...).By(op)
	}
	_ = f
	_ = ok
	b.StopTimer()
}

func Benchmark_Map_Slice(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = slice.Convert(values, op)
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Iterator(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.ToSlice(it.Convert(it.Of(values...), op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Iterator_Impl(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterimpl.ToSlice[string](iterimpl.Convert(iterimpl.New(values), op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.ToSlice(op.Convert(items, concat))
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator_No_Cache_Operation(b *testing.B) {
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.ToSlice(op.Convert(items, convert.And(toString, addTail)))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator_Impl(b *testing.B) {
	op := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterimpl.ToSlice[string](iterimpl.Convert(ptr.Of(items.Head()), op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Vector_ForEach_(b *testing.B) {
	op := convert.And(toString, addTail)
	items := vector.Of(values...)
	c := len(values)
	var s *mutable.Vector[string]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = mvector.NewCap[string](c)
		items.ForEach(func(element int) { s.Add(op(element)) })
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
		s = it.ToSlice(it.Convert(it.Filter(it.Wrap(items), even), convert.And(toString, addTail)))
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
		s = iterimpl.ToSlice[string](iterimpl.Convert(iterimpl.Filter(iterimpl.New(items), even), convert.And(toString, addTail)))
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
		s = it.ToSlice(it.Convert(sliceit.Filter(items, even), convert.And(toString, addTail)))
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
		f := sliceitimpl.Filter(items, even)
		s = iterimpl.ToSlice[string](iterimpl.Convert(&f, convert.And(toString, addTail)))
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
				s = append(s, convert.And(toString, addTail)(i))
			}
		}
	}
	_ = s
	// fmt.Println(b.Name(), s)
	b.StopTimer()
}

func Benchmark_FilterAndConvert_Iterable(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.ToSlice(it.FilterAndConvert(it.Wrap(items), even, convert.And(toString, addTail)))

	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_FilterAndConvert_Iterable_Impl(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iterimpl.NewHead(items)
		s = iterimpl.ToSlice[string](iterimpl.FilterAndConvert(&it, even, convert.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_FilterAndConvert_Embedder_Slice(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = slice.FilterAndConvert(items, even, convert.And(toString, addTail))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_FilterAndConvert_Slice(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = it.ToSlice(sliceit.FilterAndConvert(items, even, convert.And(toString, addTail)))
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_FilterAndConvert_Slice_Impl(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := sliceitimpl.FilterAndConvert(items, even, convert.And(toString, addTail))
		s = it.ToSlice[string](&m)
	}
	_ = s

	// fmt.Println(b.Name(), s)

	b.StopTimer()
}

func Benchmark_Flatt_Iterable(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := it.ToSlice(it.Filter(it.Flatt(it.Flatt(it.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Iterable_Impl(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iterimpl.NewHead(multiDimension)
		twoD := iterimpl.Flatt(&it, convert.To[[][]int])
		oneD := iterimpl.Flatt(&twoD, convert.To[[]int])
		oneDimension := iterimpl.ToSlice[int](iterimpl.Filter(&oneD, odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := it.ToSlice(it.Filter(it.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_Impl(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneD := iterimpl.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int])
		oneDimension := iterimpl.ToSlice[int](iterimpl.Filter(&oneD, odds))
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
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		result = it.Reduce(it.Filter(it.Flatt(it.Flatt(it.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds), sop.Sum[int])
	}
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Iterable_Impl(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		it := iterimpl.NewHead(multiDimension)
		twoD := iterimpl.Flatt(&it, convert.To[[][]int])
		oneD := iterimpl.Flatt(&twoD, convert.To[[]int])
		result = iterimpl.Reduce(iterimpl.Filter(&oneD, odds), sop.Sum[int])
	}
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		result = it.Reduce(it.Filter(it.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds), sop.Sum[int])
	}
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice_Impl(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	odds := func(v int) bool { return v%2 != 0 }
	toTwoD := convert.To[[][]int]
	toOneD := convert.To[[]int]
	intSum := sop.Sum[int]
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		twoD := sliceitimpl.Flatt(multiDimension, toTwoD)
		oneD := iterimpl.Flatt(&twoD, toOneD)
		result = iterimpl.Reduce(iterimpl.Filter(&oneD, odds), intSum)
	}
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	isOdd := func(val int) bool { return val%2 != 0 }
	sum := func(a, v int) int { return a + v }
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for j := 0; j < b.N; j++ {
		result = 0
		for _, i := range multiDimension {
			for _, ii := range i {
				for _, iii := range ii {
					if isOdd(iii % 2) {
						result = sum(result, iii)
					}
				}
			}
		}
	}
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
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
		_ = it.ToSlice(it.Convert(it.NotNil[Attributes](it.Flatt(it.NotNil[Participant](it.Wrap(items)), (*Participant).GetAttributes)), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableWithoutNotNilFiltering(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.ToSlice(it.Convert(it.Flatt(it.Wrap(items), (*Participant).GetAttributes), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	expected := []string{"first", "second"}
	result := []string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = it.ToSlice(it.FilterAndConvert(it.FilterAndFlatt(it.Wrap(items), check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
	}
	if !reflect.DeepEqual(expected, result) {
		b.Fatalf("must be %v, but %v", expected, result)
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFit_Impl(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iterimpl.NewHead(items)
		attr := iterimpl.FilterAndFlatt(&it, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = iterimpl.ToSlice[string](iterimpl.FilterAndConvert(&attr, check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = it.ToSlice(it.FilterAndConvert(sliceit.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit_Impl(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceitimpl.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = iterimpl.ToSlice[string](iterimpl.FilterAndConvert(&att, check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_Slice_PlainOld(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

	flatt := func(items []*Participant) []string {
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
		return names
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_ = flatt(items)
	}
	b.StopTimer()
}
