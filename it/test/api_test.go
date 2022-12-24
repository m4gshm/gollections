package it

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/conv"
	"github.com/m4gshm/gollections/it"
	impl "github.com/m4gshm/gollections/it/impl/it"
	sliceit "github.com/m4gshm/gollections/it/slice"
	"github.com/m4gshm/gollections/kvit"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
)

func Test_MapAndFilter(t *testing.T) {

	var (
		toString = func(i int) string { return fmt.Sprintf("%d", i) }
		addTail  = func(s string) string { return s + "_tail" }
	)
	items := []int{1, 2, 3, 4, 5}
	converted := it.MapFit(it.Wrap(items), func(v int) bool { return v%2 == 0 }, conv.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), it.Slice(converted))

	converted2 := sliceit.MapFit(items, func(v int) bool { return v%2 == 0 }, conv.And(toString, addTail))
	assert.Equal(t, slice.Of("2_tail", "4_tail"), it.Slice(converted2))

	//plain old style
	convertedOld := make([]string, 0)
	for _, i := range items {
		if i%2 == 0 {
			convertedOld = append(convertedOld, conv.And(toString, addTail)(i))
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
	a := it.Slice(it.Filter(it.Flatt(it.Flatt(it.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds))
	assert.Equal(t, expected, a)

	a = it.Slice(it.Filter(it.Flatt(sliceit.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds))
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

	oddSum := it.Reduce(it.Filter(it.Flatt(it.Flatt(it.Wrap(multiDimension), conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
	assert.Equal(t, e, oddSum)

	oddSum = it.Reduce(it.Filter(it.Flatt(sliceit.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
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

func Test_MapFlattStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}

	names := it.Slice(it.Map(it.Flatt(it.Wrap(items), (*Participant).GetAttributes), (*Attributes).GetName))
	assert.Equal(t, expected, names)

	names = it.Slice(it.Map(sliceit.Flatt(items, (*Participant).GetAttributes), (*Attributes).GetName))
	assert.Equal(t, expected, names)
}

func Test_MapFlattFitStructure_Iterable(t *testing.T) {
	expected := slice.Of("first", "second", "", "third", "")

	items := []*Participant{{attributes: []*Attributes{{name: "first"}, {name: "second"}, nil}}, nil, {attributes: []*Attributes{{name: "third"}, nil}}}

	names := it.Slice(it.Map(it.FlattFit(it.Wrap(items), check.NotNil[Participant], (*Participant).GetAttributes), (*Attributes).GetName))
	assert.Equal(t, expected, names)

	names = it.Slice(it.Map(sliceit.FlattFit(items, check.NotNil[Participant], (*Participant).GetAttributes), (*Attributes).GetName))
	assert.Equal(t, expected, names)
}

func Test_Iterate(t *testing.T) {
	amount := 100
	values := make([]int, amount)
	for i := 0; i < amount; i++ {
		values[i] = i
	}

	stream := impl.NewPipe[int](ptr.Of(impl.NewHead(values)))

	result := make([]int, 0)

	stream.ForEach(func(i int) { result = append(result, i) })

	result = make([]int, 0)
	it.ForEach(it.Wrap(values), func(i int) { result = append(result, i) })

	assert.Equal(t, values, result)

}

func Test_Group(t *testing.T) {
	groups := it.Group(it.Of(0, 1, 1, 2, 4, 3, 1, 6, 7), func(e int) bool { return e%2 == 0 }).Collect()

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []int{1, 1, 3, 1, 7}, groups[false])
	assert.Equal(t, []int{0, 2, 4, 6}, groups[true])
}

func Test_ReduceSum(t *testing.T) {
	s := it.Of(1, 3, 5, 7, 9, 11)
	r := it.Reduce(s, op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_First(t *testing.T) {
	s := it.Of(1, 3, 5, 7, 9, 11)
	r, ok := it.First(s, func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := it.First(s, func(i int) bool { return i > 12 })
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
	iter := it.OfLoop(stream, (*rows[int]).hasNext, (*rows[int]).next)
	s := it.Slice[int](iter)

	assert.Equal(t, slice.Of(1, 2, 3), s)
	assert.Nil(t, iter.Error())

	streamWithError := &rows[int]{slice.Of(1, 2, 3, 4), 0}
	iterWithError := it.OfLoop(streamWithError, (*rows[int]).hasNext, (*rows[int]).next)
	s2 := it.Slice[int](iterWithError)

	assert.Equal(t, slice.Of(1, 2, 3), s2)
	assert.Equal(t, "next error", iterWithError.Error().Error())
}

func Test_ToKVIterConvert(t *testing.T) {
	var (
		iter   = it.Of(K.V("key1", 1), K.V("key2", 2))
		kvIter = it.ToKVIter[string, int](iter)
		m      = kvit.Map(kvIter)
	)

	assert.Equal(t, 1, m["key1"])
	assert.Equal(t, 2, m["key2"])

}
