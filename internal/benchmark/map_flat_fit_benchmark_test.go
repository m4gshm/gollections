package benchmark

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/loop"
	sop "github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/slice"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
)

var (
	toString = func(i int) string { return fmt.Sprintf("%d", i) }
	addTail  = func(s string) string { return s + "_tail" }
	even     = func(v int) bool { return v%2 == 0 }
)

func Benchmark_First_PlainOld(b *testing.B) {
	op := func(i int) bool { return i > threshhold }
	var f int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			if op(v) {
				f = v
				break
			}
		}
	}
	b.StopTimer()
	assert.Equal(b, threshhold+1, f)
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

func Benchmark_ConvertAndFilter_Iterable(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = loop.Slice(iter.Convert(iter.Filter(slice.NewIter(items), even), convert.And(toString, addTail)).Next)
	}
	_ = s

	b.StopTimer()
}

func Benchmark_ConvertAndFilter_Slice_Loop(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = loop.Slice(loop.Convert(loop.Filter(slice.NewIter(items).Next, even).Next, convert.And(toString, addTail)).Next)
	}
	_ = s

	b.StopTimer()
}

func Benchmark_ConvertAndFilter_Slice_Iterated(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ft := sliceIter.Filter(items, even)
		conv := iter.Convert(ft, convert.And(toString, addTail))
		s = loop.Slice(conv.Next)
	}
	_ = s

	b.StopTimer()
}

func Benchmark_ConvertAndFilter_Slice_Impl(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ft := sliceIter.Filter(items, even)
		s = loop.Slice(loop.Convert(ft.Next, convert.And(toString, addTail)).Next)
	}
	_ = s

	b.StopTimer()
}

func Benchmark_ConvertAndFilter_Slice_PlainOld(b *testing.B) {
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
		s = loop.Slice[string](iter.FilterAndConvert(slice.NewIter(items), even, convert.And(toString, addTail)).Next)

	}
	_ = s

	b.StopTimer()
}

func Benchmark_FilterAndConvert_Loop(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := loop.Of(items...)
		s = loop.Slice(loop.FilterAndConvert(next, even, convert.And(toString, addTail)).Next)
	}
	_ = s

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
		f := sliceIter.FilterAndConvert(items, even, convert.And(toString, addTail))
		s = loop.Slice[string](f.Next)
	}
	_ = s

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
		m := sliceIter.FilterAndConvert(items, even, convert.And(toString, addTail))
		s = loop.Slice(m.Next)
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Flatt_Iterable(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := loop.Slice[int](iter.Filter(iter.Flat(iter.Flat(slice.NewIter(multiDimension), as.Is[[][]int]), as.Is), odds).Next)
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Loop(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := loop.Of(multiDimension...)
		twoD := loop.Flat(next, as.Is[[][]int])
		oneD := loop.Flat(twoD.Next, as.Is)
		f := loop.Filter(oneD.Next, odds)
		oneDimension := loop.Slice(f.Next)
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_Iterated(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := loop.Slice[int](iter.Filter(iter.Flat(sliceIter.Flat(multiDimension, as.Is[[][]int]), as.Is), odds).Next)
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_Looped(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf := sliceIter.Flat(multiDimension, as.Is[[][]int])
		oneD := loop.Flat(sf.Next, as.Is)
		f := loop.Filter(oneD.Next, odds)
		oneDimension := loop.Slice(f.Next)
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
		result = loop.Reduce(iter.Filter(iter.Flat(iter.Flat(slice.NewIter(multiDimension), as.Is[[][]int]), as.Is), odds).Next, sop.Sum[int])
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
}

func Benchmark_ReduceSum_Loop(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		next := loop.Of(multiDimension...)
		f1 := loop.Flat(next, as.Is[[][]int])
		f2 := loop.Flat(f1.Next, as.Is)
		f3 := loop.Filter(f2.Next, odds)
		result = loop.Reduce(f3.Next, sop.Sum[int])
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
}

func Benchmark_ReduceSum_Slice(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		result = loop.Reduce(iter.Filter(iter.Flat(sliceIter.Flat(multiDimension, as.Is[[][]int]), as.Is), odds).Next, sop.Sum[int])
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
}

//go:noinline
func odds(v int) bool { return v%2 != 0 }

func Benchmark_ReduceSum_Slice_Looped(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7

	b.ResetTimer()
	result := 0
	for i := 0; i < b.N; i++ {
		result = loop.Reduce(loop.Filter(loop.Flat(sliceIter.Flat(multiDimension, as.Is).Next, as.Is).Next, odds).Next, sop.Sum)
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
}

func Benchmark_ReduceSum_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for j := 0; j < b.N; j++ {
		result = 0
		for _, i := range multiDimension {
			for _, ii := range i {
				for _, iii := range ii {
					if odds(iii % 2) {
						result = sop.Sum(result, iii)
					}
				}
			}
		}
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
}

func Benchmark_ReduceSum_Slice_PlainOld_Index(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for j := 0; j < b.N; j++ {
		result = 0
		for i := range multiDimension {
			for ii := range multiDimension[i] {
				for iii := range multiDimension[i][ii] {
					if odds(multiDimension[i][ii][iii] % 2) {
						result = sop.Sum(result, multiDimension[i][ii][iii])
					}
				}
			}
		}
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
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

func Benchmark_ConvertFlattStructure_IterableNotNil(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loop.Slice(iter.Convert(iter.NotNil[Attributes](iter.Flat(iter.NotNil[Participant](slice.NewIter(items)), (*Participant).GetAttributes)), (*Attributes).GetName).Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_IterableWithoutNotNilFiltering(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = loop.Slice(iter.Convert(iter.Flat(slice.NewIter(items), (*Participant).GetAttributes), (*Attributes).GetName).Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_IterableFilt(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	expected := []string{"first", "second"}
	result := []string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = loop.Slice(iter.FilterAndConvert(iter.FilterAndFlat(slice.NewIter(items), not.Nil[Participant], (*Participant).GetAttributes), not.Nil[Attributes], (*Attributes).GetName).Next)
	}
	if !reflect.DeepEqual(expected, result) {
		b.Fatalf("must be %v, but %v", expected, result)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_Loop(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		next := loop.Of(items...)
		attr := loop.FilterAndFlat(next, not.Nil[Participant], (*Participant).GetAttributes)
		f := loop.FilterAndConvert(attr.Next, not.Nil[Attributes], (*Attributes).GetName)
		_ = loop.Slice(f.Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_SliceFilt(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceIter.FilterAndFlat(items, not.Nil[Participant], (*Participant).GetAttributes)
		_ = loop.Slice(iter.FilterAndConvert(att, not.Nil[Attributes], (*Attributes).GetName).Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_SliceFilt_Looped(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceIter.FilterAndFlat(items, not.Nil[Participant], (*Participant).GetAttributes)
		_ = loop.Slice(loop.FilterAndConvert(att.Next, not.Nil[Attributes], (*Attributes).GetName).Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_Slice_PlainOld(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}

	flattener := func(items []*Participant) []string {
		names := make([]string, 0)
		for _, i := range items {
			if not.Nil(i) {
				for _, a := range i.GetAttributes() {
					if not.Nil(a) {
						names = append(names, a.GetName())
					}
				}
			}
		}
		return names
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_ = flattener(items)
	}
	b.StopTimer()
}
