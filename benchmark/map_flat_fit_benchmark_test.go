package benchmark

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/iter"
	iterimpl "github.com/m4gshm/gollections/iter/impl/iter"
	sliceitimpl "github.com/m4gshm/gollections/iter/impl/slice"
	sliceit "github.com/m4gshm/gollections/iter/slice"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/mutable"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	sop "github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
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
	op := func(i int) bool { return i > threshhold }
	var f int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = slice.First(values, op)
	}
	b.StopTimer()
	assert.Equal(b, threshhold+1, f)
}

func Benchmark_First_OfBy(b *testing.B) {
	op := func(i int) bool { return i > threshhold }
	var f int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = first.Of(values...).By(op)
	}

	b.StopTimer()
	assert.Equal(b, threshhold+1, f)
}

func Benchmark_First_Iterator(b *testing.B) {
	op := func(i int) bool { return i > threshhold }
	var f int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = iter.First(iter.Of(values...), op)
	}
	b.StopTimer()
	assert.Equal(b, threshhold+1, f)
}

func Benchmark_Last_Slice(b *testing.B) {
	op := func(i int) bool { return i < 50000 }
	var f int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = slice.Last(values, op)
	}
	b.StopTimer()
	assert.Equal(b, threshhold-1, f)
}

func Benchmark_Last_OfBy(b *testing.B) {
	op := func(i int) bool { return i < threshhold }
	var f int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, _ = last.Of(values...).By(op)
	}
	b.StopTimer()
	assert.Equal(b, threshhold-1, f)
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
		s = iter.ToSlice(iter.Convert(iter.Of(values...), op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Map_Iterator_Impl(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iterimpl.New(values)
		s = loop.ToSlice(iterimpl.Convert(it, it.Next, op).Next)
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
		s = iterable.Convert(items, concat).Slice()
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Map_Vector_Iterator_No_Cache_Operation(b *testing.B) {
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterable.Convert(items, convert.And(toString, addTail)).Slice()
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
		h := ptr.Of(items.Head())
		s = loop.ToSlice(iterimpl.Convert(h, h.Next, op).Next)
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
		s = iter.ToSlice(iter.Convert(iter.Filter(iter.Wrap(items), even), convert.And(toString, addTail)))
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
		it := iterimpl.New(items)
		ft := iterimpl.Filter(it, it.Next, even)
		s = loop.ToSlice(iterimpl.Convert(ft, ft.Next, convert.And(toString, addTail)).Next)
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
		s = iter.ToSlice(iter.Convert(sliceit.Filter(items, even), convert.And(toString, addTail)))
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
		ft := sliceitimpl.Filter(items, even)
		s = loop.ToSlice(iterimpl.Convert(ft, ft.Next, convert.And(toString, addTail)).Next)
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
		s = iter.ToSlice(iter.FilterAndConvert(iter.Wrap(items), even, convert.And(toString, addTail)))

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
		it := ptr.Of(iterimpl.NewHead(items))
		s = loop.ToSlice(iterimpl.FilterAndConvert(it, it.Next, even, convert.And(toString, addTail)).Next)
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
		s = iter.ToSlice(sliceit.FilterAndConvert(items, even, convert.And(toString, addTail)))
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
		s = loop.ToSlice(m.Next)
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
		oneDimension := iter.ToSlice(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds))
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
		twoD := iterimpl.Flatt(it, it.Next, convert.To[[][]int])
		oneD := iterimpl.Flatt(twoD, twoD.Next, convert.To[[]int])
		oneDimension := loop.ToSlice(iterimpl.Filter(oneD, oneD.Next, odds).Next)
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := iter.ToSlice(iter.Filter(iter.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_Impl(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf := sliceit.Flatt(multiDimension, convert.To[[][]int])
		oneD := iterimpl.Flatt(sf, sf.Next, convert.To[[]int])
		oneDimension := loop.ToSlice(iterimpl.Filter(oneD, oneD.Next, odds).Next)
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
		result = iter.Reduce(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds), sop.Sum[int])
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
		twoD := iterimpl.Flatt(it, it.Next, convert.To[[][]int])
		oneD := iterimpl.Flatt(twoD, twoD.Next, convert.To[[]int])
		result = loop.Reduce(iterimpl.Filter(oneD, oneD.Next, odds).Next, sop.Sum[int])
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
		result = iter.Reduce(iter.Filter(iter.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds), sop.Sum[int])
	}
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
	b.StopTimer()
}

func Benchmark_ReduceSum_Slice_Impl(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7

	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		odds := func(v int) bool { return v%2 != 0 }
		f1 := sliceitimpl.Flatt(multiDimension, convert.To[[][]int])
		f2 := iterimpl.Flatt(f1, f1.Next, convert.To[[]int])
		result = loop.Reduce(iterimpl.Filter(f2, f2.Next, odds).Next, sop.Sum[int])
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
		_ = iter.ToSlice(iter.Convert(iter.NotNil[Attributes](iter.Flatt(iter.NotNil[Participant](iter.Wrap(items)), (*Participant).GetAttributes)), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableWithoutNotNilFiltering(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice(iter.Convert(iter.Flatt(iter.Wrap(items), (*Participant).GetAttributes), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_IterableFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	expected := []string{"first", "second"}
	result := []string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = iter.ToSlice(iter.FilterAndConvert(iter.FilterAndFlatt(iter.Wrap(items), check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
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
		attr := iterimpl.FilterAndFlatt(&it, (&it).Next, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = loop.ToSlice(iterimpl.FilterAndConvert(&attr, (&attr).Next, check.NotNil[Attributes], (*Attributes).GetName).Next)
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceit.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = iter.ToSlice(iter.FilterAndConvert(att, check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_MapFlattStructure_SliceFit_Impl(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceitimpl.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = loop.ToSlice(iterimpl.FilterAndConvert(att, att.Next, check.NotNil[Attributes], (*Attributes).GetName).Next)
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
