package test

import (
	"testing"

	"github.com/m4gshm/gollections/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/iter"
	"github.com/stretchr/testify/assert"
)

func Test_PointerBasedIter(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	expected := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	iterator := iter.NewHead(expected)
	result := make([]someType, 0)
	for iterator.HasNext() {
		n := iterator.GetNext()
		result = append(result, n)
	}

	assert.Equal(t, expected, result)
}

func Test_PointerBasedIter2(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	expected := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	iterator := iter.NewHead(expected)
	result := make([]someType, 0)
	for v, ok := iterator.Next(); ok; v, ok = iterator.Next() {
		result = append(result, v)
	}

	assert.Equal(t, expected, result)
}

func Test_PointerBasedIter2Reverse(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	expected := []someType{{"4", 4}, {"3", 3}, {"2", 2}, {"123", 123}}
	iterator := iter.NewTail(values)
	result := make([]someType, 0)
	for v, ok := iterator.Prev(); ok; v, ok = iterator.Prev() {
		result = append(result, v)
	}

	assert.Equal(t, expected, result)
}

func Test_PointerBasedIterTailGetNext(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	iterator := iter.NewTail(values)

	v, ok := iterator.Get() //out of range
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Next() //no more elements, because the iterator has not been started
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Prev() //start from the latest element
	assert.True(t, ok)
	assert.Equal(t, someType{"4", 4}, v)

	v, ok = iterator.Get() //gets the current element
	assert.True(t, ok)
	assert.Equal(t, someType{"4", 4}, v)

	v, ok = iterator.Next()
	_ = v
	assert.False(t, ok) //no more elements after the latest

	v, ok = iterator.Prev() //gets prev
	assert.Equal(t, someType{"3", 3}, v)
	assert.True(t, ok)

	v, ok = iterator.Get() //gets the current element
	assert.True(t, ok)
	assert.Equal(t, someType{"3", 3}, v)

	v, ok = iterator.Next()
	assert.True(t, ok) //returns to the latest
	assert.Equal(t, someType{"4", 4}, v)

	v, ok = iterator.Get() //gets the current element
	assert.True(t, ok)
	assert.Equal(t, someType{"4", 4}, v)
}

func Test_PointerBasedIterHeadGetPrev(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	iterator := iter.NewHead(values)

	v, ok := iterator.Prev()
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Next()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)

	v, ok = iterator.Get()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)

	v, ok = iterator.Prev()
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Next()
	assert.True(t, ok)
	assert.Equal(t, someType{"2", 2}, v)

	v, ok = iterator.Get()
	assert.True(t, ok)
	assert.Equal(t, someType{"2", 2}, v)

	v, ok = iterator.Prev()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)

	v, ok = iterator.Get()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)
}

func Test_PointerBasedEmptyIter(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{}

	iterator := iter.NewHead(values)

	v, ok := iterator.Prev()
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Next()
	assert.False(t, ok)

	v, ok = iterator.Get()
	assert.False(t, ok)

	//tail

	iterator = iter.NewTail(values)

	v, ok = iterator.Prev()
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Next()
	assert.False(t, ok)

	v, ok = iterator.Get()
	assert.False(t, ok)
}

func Test_PointerBasedOneElementIter(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{{"only one", 1}}
	iterator := iter.NewHead(values)

	v, ok := iterator.Prev()
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Next()
	assert.True(t, ok)

	v, ok = iterator.Get()
	assert.True(t, ok)

	iterator = iter.NewTail(values)

	v, ok = iterator.Next()
	_ = v
	assert.False(t, ok)

	v, ok = iterator.Prev()
	assert.True(t, ok)

	v, ok = iterator.Get()
	assert.True(t, ok)
}

func Test_CanIterateByRange(t *testing.T) {
	r := iter.CanIterateByRange(iter.NoStarted, 5, 4)
	assert.True(t, r)

	r = iter.CanIterateByRange(iter.NoStarted, 5, 6)
	assert.False(t, r)

	r = iter.CanIterateByRange(iter.NoStarted, 5, iter.NoStarted)
	assert.True(t, r)
}

func Test_IsValidIndex(t *testing.T) {
	r := iter.IsValidIndex(5, 0)
	assert.True(t, r)

	r = iter.IsValidIndex(5, 5)
	assert.False(t, r)

	r = iter.IsValidIndex(5, -1)
	assert.False(t, r)

}

func Test_IsValidIndex2(t *testing.T) {
	r := iter.IsValidIndex2(5, 0)
	assert.True(t, r)

	r = iter.IsValidIndex2(5, 5)
	assert.False(t, r)

	r = iter.IsValidIndex2(5, -1)
	assert.False(t, r)
}

func Test_Head_Tail_Nil_Arg_Safety(t *testing.T) {
	var values []int

	head := iter.NewHead(values)

	assert.False(t, head.HasNext())
	assert.False(t, head.HasPrev())

	_, ok := head.Get()
	assert.False(t, ok)

	_, ok = head.Next()
	assert.False(t, ok)
	_, ok = head.Prev()
	assert.False(t, ok)
	head.Cap()

	tail := iter.NewTail(values)

	assert.False(t, tail.HasNext())
	assert.False(t, tail.HasPrev())
	_, ok = tail.Get()
	assert.False(t, ok)
	_, ok = tail.Next()
	assert.False(t, ok)
	_, ok = tail.Prev()
	assert.False(t, ok)
	tail.Cap()
}

func Test_ForLoop(t *testing.T) {

	expected := []int{1, 2, 3, 4, 5, 6, 7}
	actual := []int{}
	it := iter.NewHead(expected)
	it.For(func(element int) error { actual = append(actual, element); return nil })
	assert.Equal(t, expected, actual)

	actual2 := []int{}
	it = iter.NewHead(expected)
	it.ForEach(func(element int) { actual2 = append(actual2, element) })
	assert.Equal(t, expected, actual2)
}

func Test_group_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		it     = iter.NewKeyValuer(slice.Of(1, 1, 2, 4, 3, 1), even, as.Is[int])
		groups = kvloop.Group(it.Next)
	)
	assert.Equal(t, map[bool][]int{false: {1, 1, 3, 1}, true: {2, 4}}, groups)
}
