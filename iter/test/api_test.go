package it

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/break/predicate"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/slice"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
	"github.com/m4gshm/gollections/stream"
)

func Test_FilterAndConvert(t *testing.T) {

	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
	)
	items := []int{1, 2, 3, 4, 5}
	converted := iter.FilterAndConvert(slice.NewIter(items), func(v int) bool { return v%2 == 0 }, convert.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), loop.Slice(converted.Next))

	converted2 := sliceIter.FilterAndConvert(items, func(v int) bool { return v%2 == 0 }, convert.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), loop.Slice(converted2.Next))

	//plain old style
	convertedOld := make([]string, 0)
	for _, i := range items {
		if i%2 == 0 {
			convertedOld = append(convertedOld, convert.And(toString, addTail)(i))
		}
	}

	assert.Equal(t, slice.Of("2_tail", "4_tail"), convertedOld)
}

func Test_FiltAndConv(t *testing.T) {
	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
	)
	items := []int{1, 2, 3, 4, 5}

	converted := sliceIter.FiltAndConv(items, func(v int) (bool, error) { return v%2 == 0, nil }, wrap(convert.And(toString, addTail)))
	s, _ := breakLoop.Slice(converted.Next)
	assert.Equal(t, slice.Of("2_tail", "4_tail"), s)
}

func wrap[F, T any](f func(F) T) func(F) (T, error) {
	return func(i F) (T, error) { return f(i), nil }
}

func Test_FlattSlices(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
		expected       = slice.Of(1, 3, 5, 7)
	)
	f := iter.Filter(iter.Flat(iter.Flat(slice.NewIter(multiDimension), as.Is[[][]int]), as.Is), odds)
	a := loop.Slice(f.Next)
	assert.Equal(t, expected, a)

	a = loop.Slice(iter.Filter(iter.Flat(sliceIter.Flat(multiDimension, as.Is[[][]int]), as.Is), odds).Next)
	assert.Equal(t, expected, a)

	//plain old style
	oneDimensionOld := make([]int, 0)
	for _, i := range multiDimension {
		if i == nil {
			continue
		}
		for _, ii := range i {
			if ii == nil {
				continue
			}
			for _, iii := range ii {
				if odds(iii) {
					oneDimensionOld = append(oneDimensionOld, iii)
				}
			}
		}
	}

	assert.Equal(t, slice.Of(1, 3, 5, 7), oneDimensionOld)

}

func Test_ReduceSlices(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	e := 1 + 3 + 5 + 7

	oddSum := loop.Reduce(iter.Filter(iter.Flat(iter.Flat(slice.NewIter(multiDimension), as.Is[[][]int]), as.Is), odds).Next, op.Sum[int])
	assert.Equal(t, e, oddSum)

	oddSum = loop.Reduce(iter.Filter(iter.Flat(sliceIter.Flat(multiDimension, as.Is[[][]int]), as.Is), odds).Next, op.Sum[int])
	assert.Equal(t, e, oddSum)

	//plain old style
	oddSum = 0
	for _, i := range multiDimension {
		for _, ii := range i {
			for _, iii := range ii {
				if odds(iii) {
					oddSum += iii
				}
			}
		}
	}

	assert.Equal(t, e, oddSum)

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

func Test_ConvertFlattStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}

	names := loop.Slice(iter.Convert(iter.Flat(slice.NewIter(items), (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)

	names = loop.Slice(iter.Convert(sliceIter.Flat(items, (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)
}

func Test_ConvertFlatStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}

	names, _ := breakLoop.Slice(breakLoop.Convert(sliceIter.Flatt(items, wrapGet((*Participant).GetAttributes)).Next, (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)
}

func Test_ConvertFilterAndFlattStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}

	names := loop.Slice(iter.Convert(iter.FilterAndFlat(slice.NewIter(items), not.Nil[Participant], (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)

	names = loop.Slice(iter.Convert(sliceIter.FilterAndFlat(items, not.Nil[Participant], (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)
}

func Test_ConvertFiltAndFlattStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}
	names, _ := breakLoop.Slice(breakLoop.Convert(sliceIter.FiltAndFlat(items, predicate.Wrap(not.Nil[Participant]), wrapGet((*Participant).GetAttributes)).Next, (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)
}

func wrapGet[S, V any](getter func(S) V) func(S) (V, error) {
	return func(s S) (V, error) {
		return getter(s), nil
	}
}

func Test_Iterate(t *testing.T) {
	amount := 100
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}

	stream := stream.New(sliceIter.New(values).Next)

	result := make([]int, 0)

	stream.ForEach(func(i int) { result = append(result, i) })

	result = make([]int, 0)
	sliceIter.New(values).ForEach(func(i int) { result = append(result, i) })

	assert.Equal(t, values, result)

}

func Test_Group(t *testing.T) {
	groups := iter.Group(iter.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 }).Map()

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 1, 3, 1, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_ReduceSum(t *testing.T) {
	s := iter.Of(1, 3, 5, 7, 9, 11)
	r := loop.Reduce(s.Next, op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_First(t *testing.T) {
	s := iter.Of(1, 3, 5, 7, 9, 11)
	r, ok := loop.First(s.Next, func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := loop.First(s.Next, func(i int) bool { return i > 12 })
	assert.False(t, nook)
}

type rows[T any] struct {
	in     []T
	cursor int
}

func (r *rows[T]) hasNext() bool {
	return r.cursor < len(r.in)
}

func (r *rows[T]) next() (T, error) {
	e := r.in[r.cursor]
	r.cursor++
	if r.cursor > 3 {
		var no T
		return no, errors.New("next error")
	}
	return e, nil
}

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	iter := breakLoop.New(stream, (*rows[int]).hasNext, (*rows[int]).next)
	s, err := breakLoop.Slice(iter)

	assert.Equal(t, slice.Of(1, 2, 3), s)
	assert.Nil(t, err)

	streamWithError := &rows[int]{slice.Of(1, 2, 3, 4), 0}
	iterWithError := breakLoop.New(streamWithError, (*rows[int]).hasNext, (*rows[int]).next)
	s2, err2 := breakLoop.Slice(iterWithError)

	assert.Equal(t, slice.Of(1, 2, 3), s2)
	assert.Equal(t, "next error", err2.Error())
}
