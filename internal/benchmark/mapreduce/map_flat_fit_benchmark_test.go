package mapreduce

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	sop "github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/range_"
)

var (
	toString      = func(i int) string { return fmt.Sprintf("%d", i) }
	addTail       = func(s string) string { return s + "_tail" }
	even          = func(v int) bool { return v%2 == 0 }
	maxValOfRange = 100000
	values        = range_.Closed(1, maxValOfRange)
	threshold     = maxValOfRange / 2
)

func Benchmark_First_PlainOld(b *testing.B) {
	op := func(i int) bool { return i > threshold }
	var f int

	for b.Loop() {
		for _, v := range values {
			if op(v) {
				f = v
				break
			}
		}
	}
	b.StopTimer()
	assert.Equal(b, threshold+1, f)
}

func Benchmark_First_Slice(b *testing.B) {
	op := func(i int) bool { return i > threshold }
	var f int

	for b.Loop() {
		f, _ = slice.First(values, op)
	}
	b.StopTimer()
	assert.Equal(b, threshold+1, f)
}

func Benchmark_FirstI_Slice(b *testing.B) {
	op := func(i int) bool { return i > threshold }
	var f, ind int

	for b.Loop() {
		f, ind = slice.FirstI(values, op)
	}
	b.StopTimer()
	assert.Equal(b, threshold+1, f)
	assert.Equal(b, threshold, ind)
}

func Benchmark_Last_Slice(b *testing.B) {
	op := func(i int) bool { return i < 50000 }
	var f int

	for b.Loop() {
		f, _ = slice.Last(values, op)
	}
	b.StopTimer()
	assert.Equal(b, threshold-1, f)
}

func Benchmark_LastI_Slice(b *testing.B) {
	op := func(i int) bool { return i < 50000 }
	var f, ind int

	for b.Loop() {
		f, ind = slice.LastI(values, op)
	}
	b.StopTimer()
	assert.Equal(b, threshold-1, f)
	assert.Equal(b, threshold-2, ind)
}

func Benchmark_ConvertAndFilter_Slice_Seq(b *testing.B) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
		even     = func(v int) bool { return v%2 == 0 }
	)
	items := slice.Of(1, 2, 3, 4, 5)
	var s []string

	for b.Loop() {
		s = seq.Slice(seq.Convert(seq.Filter(seq.Of(items...), even), convert.And(toString, addTail)))
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
	for b.Loop() {
		s := make([]string, 0)
		for _, i := range items {
			if i%2 == 0 {
				s = append(s, convert.And(toString, addTail)(i))
			}
		}
		_ = s
	}
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
	for b.Loop() {
		s = slice.FilterAndConvert(items, even, convert.And(toString, addTail))
	}
	_ = s

	b.StopTimer()
}

func Benchmark_Flatt_Seq(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	b.ResetTimer()
	for b.Loop() {
		next := seq.Of(multiDimension...)
		oneDimension := seq.Slice(seq.Filter(seq.Flat(seq.Flat(next, as.Is), as.Is), odds))
		_ = oneDimension
	}
	b.StopTimer()
}

func Benchmark_Flatt_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}

	for b.Loop() {
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
		_ = oneDimension
	}
	b.StopTimer()
}

//go:noinline
func odds(v int) bool { return v%2 != 0 }

func Benchmark_ReduceSum_Seq(b *testing.B) {
	odds := func(v int) bool { return v%2 != 0 }
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7
	b.ResetTimer()
	result := 0
	for b.Loop() {
		result = seq.Flat(seq.Flat(seq.Of(multiDimension...), as.Is), as.Is).Filter(odds).Reduce(sop.Sum)
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
	for b.Loop() {
		result = slice.Reduce(slice.Filter(slice.Flat(slice.Flat(multiDimension, as.Is), as.Is), odds), sop.Sum)
	}
	b.StopTimer()
	if result != expected {
		b.Fatalf("must be %d, but %d", expected, result)
	}
}

func Benchmark_ReduceSum_Slice_PlainOld(b *testing.B) {
	multiDimension := [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	expected := 1 + 3 + 5 + 7

	result := 0
	for b.Loop() {
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

	result := 0
	for b.Loop() {
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

func Benchmark_ConvertFlattStructure_Seq(b *testing.B) {
	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil}
	b.ResetTimer()
	for b.Loop() {
		attr := seq.Flat(seq.Filter(seq.Of(items...), not.Nil), (*Participant).GetAttributes)
		_ = seq.Slice(seq.Convert(seq.Filter(attr, not.Nil), (*Attributes).GetName))
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

	for b.Loop() {
		_ = flattener(items)
	}
	b.StopTimer()
}
