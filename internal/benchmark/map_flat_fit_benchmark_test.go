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
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/mutable"
	mvector "github.com/m4gshm/gollections/mutable/vector"
	sop "github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
	"github.com/stretchr/testify/assert"
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

func Benchmark_Convert_Slice(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = slice.Convert(values, op)
	}
	_ = s
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

func Benchmark_Convert_Iterator(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iter.ToSlice[string](iter.Convert(iter.Of(values...), op))
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_SliceIter(b *testing.B) {
	op := convert.And(toString, addTail)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = loop.ToSlice(sliceIter.Convert(values, op).Next)
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
		s = loop.ToSlice(loop.Convert(it.Next, op).Next)
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_Vector_Iterator(b *testing.B) {
	concat := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterable.Convert[*slice.Iter[int]](items, concat).Slice()
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Convert_Vector_Iterator_No_Cache_Operation(b *testing.B) {
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = iterable.Convert[*slice.Iter[int]](items, convert.And(toString, addTail)).Slice()
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_Vector_Loop(b *testing.B) {
	op := convert.And(toString, addTail)
	items := vector.Of(values...)
	var s []string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := items.Head()
		s = loop.ToSlice(loop.Convert(h.Next, op).Next)
	}
	_ = s
	b.StopTimer()
}

func Benchmark_Convert_Vector_ForEach_(b *testing.B) {
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
		s = iter.ToSlice[string](iter.Convert(iter.Filter(iter.Wrap(items), even), convert.And(toString, addTail)))
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
		it := slice.NewIter(items)
		ft := loop.Filter(it.Next, even)
		s = loop.ToSlice(loop.Convert(ft.Next, convert.And(toString, addTail)).Next)
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
		s = iter.ToSlice[string](iter.Convert(sliceIter.Filter(items, even), convert.And(toString, addTail)))
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
		s = loop.ToSlice(loop.Convert(ft.Next, convert.And(toString, addTail)).Next)
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
		s = iter.ToSlice[string](iter.FilterAndConvert(iter.Wrap(items), even, convert.And(toString, addTail)))

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
		s = loop.ToSlice(loop.FilterAndConvert(next, even, convert.And(toString, addTail)).Next)
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
		s = iter.ToSlice[string](f)
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
		s = loop.ToSlice(m.Next)
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Flatt_Iterable(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := iter.ToSlice[int](iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds))
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
		twoD := loop.Flatt(next, convert.To[[][]int])
		oneD := loop.Flatt(twoD.Next, convert.To[[]int])
		f := loop.Filter(oneD.Next, odds)
		oneDimension := loop.ToSlice(f.Next)
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_Iterated(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		oneDimension := iter.ToSlice[int](iter.Filter(iter.Flatt(sliceIter.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_Looped(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf := sliceIter.Flatt(multiDimension, convert.To[[][]int])
		oneD := loop.Flatt(sf.Next, convert.To[[]int])
		f := loop.Filter(oneD.Next, odds)
		oneDimension := loop.ToSlice(f.Next)
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
		result = loop.Reduce(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds).Next, sop.Sum[int])
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
		f1 := loop.Flatt(next, convert.To[[][]int])
		f2 := loop.Flatt(f1.Next, convert.To[[]int])
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
		result = loop.Reduce(iter.Filter(iter.Flatt(sliceIter.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds).Next, sop.Sum[int])
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
		f1 := sliceIter.Flatt(multiDimension, convert.To[[][]int])
		f2 := loop.Flatt(f1.Next, convert.To[[]int])
		f3 := loop.Filter(f2.Next, odds)
		result = loop.Reduce(f3.Next, sop.Sum[int])
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
		_ = iter.ToSlice[string](iter.Convert(iter.NotNil[Attributes](iter.Flatt(iter.NotNil[Participant](iter.Wrap(items)), (*Participant).GetAttributes)), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_IterableWithoutNotNilFiltering(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = iter.ToSlice[string](iter.Convert(iter.Flatt(iter.Wrap(items), (*Participant).GetAttributes), (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_IterableFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	expected := []string{"first", "second"}
	result := []string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = iter.ToSlice[string](iter.FilterAndConvert(iter.FilterAndFlatt(iter.Wrap(items), check.NotNil[Participant], (*Participant).GetAttributes), check.NotNil[Attributes], (*Attributes).GetName))
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
		attr := loop.FilterAndFlatt(next, check.NotNil[Participant], (*Participant).GetAttributes)
		f := loop.FilterAndConvert(attr.Next, check.NotNil[Attributes], (*Attributes).GetName)
		_ = loop.ToSlice(f.Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_SliceFit(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceIter.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = iter.ToSlice[string](iter.FilterAndConvert(att, check.NotNil[Attributes], (*Attributes).GetName))
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_SliceFit_Looped(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		att := sliceIter.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes)
		_ = loop.ToSlice(loop.FilterAndConvert(att.Next, check.NotNil[Attributes], (*Attributes).GetName).Next)
	}
	b.StopTimer()
}

func Benchmark_ConvertFlattStructure_Slice_PlainOld(b *testing.B) {
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
