package it

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/iter"
	sliceit "github.com/m4gshm/gollections/iter/slice"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
)

func Test_FilterAndConvert(t *testing.T) {

	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
	)
	items := []int{1, 2, 3, 4, 5}
	converted := iter.FilterAndConvert(iter.Wrap(items), func(v int) bool { return v%2 == 0 }, convert.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), loop.ToSlice(converted.Next))

	converted2 := sliceit.FilterAndConvert(items, func(v int) bool { return v%2 == 0 }, convert.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), loop.ToSlice(converted2.Next))

	//plain old style
	convertedOld := make([]string, 0)
	for _, i := range items {
		if i%2 == 0 {
			convertedOld = append(convertedOld, convert.And(toString, addTail)(i))
		}
	}

	assert.Equal(t, slice.Of("2_tail", "4_tail"), convertedOld)

}

func Test_FlattSlices(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
		expected       = slice.Of(1, 3, 5, 7)
	)
	a := iter.ToSlice(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds))
	assert.Equal(t, expected, a)

	a = iter.ToSlice(iter.Filter(iter.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds))
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

	oddSum := loop.Reduce(iter.Filter(iter.Flatt(iter.Flatt(iter.Wrap(multiDimension), convert.To[[][]int]), convert.To[[]int]), odds).Next, op.Sum[int])
	assert.Equal(t, e, oddSum)

	oddSum = loop.Reduce(iter.Filter(iter.Flatt(sliceit.Flatt(multiDimension, convert.To[[][]int]), convert.To[[]int]), odds).Next, op.Sum[int])
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

	names := loop.ToSlice(iter.Convert(iter.Flatt(iter.Wrap(items), (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)

	names = loop.ToSlice(iter.Convert(sliceit.Flatt(items, (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)
}

func Test_ConvertFilterAndFlattStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}

	names := loop.ToSlice(iter.Convert(iter.FilterAndFlatt(iter.Wrap(items), check.NotNil[Participant], (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)

	names = loop.ToSlice(iter.Convert(sliceit.FilterAndFlatt(items, check.NotNil[Participant], (*Participant).GetAttributes), (*Attributes).GetName).Next)
	assert.Equal(t, expected, names)
}

func Test_Iterate(t *testing.T) {
	amount := 100
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}

	stream := loop.Stream(sliceIter.New(values).Next)

	result := make([]int, 0)

	stream.ForEach(func(i int) { result = append(result, i) })

	result = make([]int, 0)
	iter.ForEach(iter.New(values), func(i int) { result = append(result, i) })

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
	r, ok := iter.First(s, func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := iter.First(s, func(i int) bool { return i > 12 })
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
	s, err := breakLoop.ToSlice(iter)

	assert.Equal(t, slice.Of(1, 2, 3), s)
	assert.Nil(t, err)

	streamWithError := &rows[int]{slice.Of(1, 2, 3, 4), 0}
	iterWithError := breakLoop.New(streamWithError, (*rows[int]).hasNext, (*rows[int]).next)
	s2, err2 := breakLoop.ToSlice(iterWithError)

	assert.Equal(t, slice.Of(1, 2, 3), s2)
	assert.Equal(t, "next error", err2.Error())
}
