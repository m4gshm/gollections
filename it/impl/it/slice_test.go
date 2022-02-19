package it

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PointerBasedIter(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	expected := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	iter := NewHead(expected)
	result := make([]someType, 0)
	for iter.HasNext() {
		n := iter.Next()
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
	iter := NewHead(expected)
	result := make([]someType, 0)
	for v, ok := iter.GetNext(); ok; v, ok = iter.GetNext() {
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
	iter := NewTail(values)
	result := make([]someType, 0)
	for v, ok := iter.GetPrev(); ok; v, ok = iter.GetPrev() {
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
	iter := NewTail(values)

	v, ok := iter.Get() //out of range
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetNext() //no more elements, because the iterator has not been started
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetPrev() //start from the latest element
	assert.True(t, ok)
	assert.Equal(t, someType{"4", 4}, v)

	v, ok = iter.Get() //gets the current element
	assert.True(t, ok)
	assert.Equal(t, someType{"4", 4}, v)

	v, ok = iter.GetNext()
	_ = v
	assert.False(t, ok) //no more elements after the latest

	v, ok = iter.GetPrev() //gets prev
	assert.Equal(t, someType{"3", 3}, v)
	assert.True(t, ok)

	v, ok = iter.Get() //gets the current element
	assert.True(t, ok)
	assert.Equal(t, someType{"3", 3}, v)

	v, ok = iter.GetNext()
	assert.True(t, ok) //returns to the latest
	assert.Equal(t, someType{"4", 4}, v)

	v, ok = iter.Get() //gets the current element
	assert.True(t, ok)
	assert.Equal(t, someType{"4", 4}, v)
}

func Test_PointerBasedIterHeadGetPrev(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{{"123", 123}, {"2", 2}, {"3", 3}, {"4", 4}}
	iter := NewHead(values)

	v, ok := iter.GetPrev()
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetNext()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)

	v, ok = iter.Get()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)

	v, ok = iter.GetPrev()
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetNext()
	assert.True(t, ok)
	assert.Equal(t, someType{"2", 2}, v)

	v, ok = iter.Get()
	assert.True(t, ok)
	assert.Equal(t, someType{"2", 2}, v)

	v, ok = iter.GetPrev()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)

	v, ok = iter.Get()
	assert.True(t, ok)
	assert.Equal(t, someType{"123", 123}, v)
}

func Test_PointerBasedEmptyIter(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{}

	iter := NewHead(values)

	v, ok := iter.GetPrev()
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetNext()
	assert.False(t, ok)

	v, ok = iter.Get()
	assert.False(t, ok)

	//tail

	iter = NewTail(values)

	v, ok = iter.GetPrev()
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetNext()
	assert.False(t, ok)

	v, ok = iter.Get()
	assert.False(t, ok)
}

func Test_PointerBasedOneElementIter(t *testing.T) {

	type someType struct {
		field1 string
		field2 int64
	}

	values := []someType{{"only one", 1}}
	iter := NewHead(values)

	v, ok := iter.GetPrev()
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetNext()
	assert.True(t, ok)

	v, ok = iter.Get()
	assert.True(t, ok)

	iter = NewTail(values)

	v, ok = iter.GetNext()
	_ = v
	assert.False(t, ok)

	v, ok = iter.GetPrev()
	assert.True(t, ok)

	v, ok = iter.Get()
	assert.True(t, ok)
}

func Test_CanIterateByRange(t *testing.T) {
	r := CanIterateByRange(NoStarted, 5, 4)
	assert.True(t, r)

	r = CanIterateByRange(NoStarted, 5, 6)
	assert.False(t, r)

	r = CanIterateByRange(NoStarted, 5, NoStarted)
	assert.True(t, r)
}

func Test_IsValidIndex(t *testing.T) {
	r := IsValidIndex(5, 0)
	assert.True(t, r)

	r = IsValidIndex(5, 5)
	assert.False(t, r)

	r = IsValidIndex(5, -1)
	assert.False(t, r)

}

func Test_IsValidIndex2(t *testing.T) {
	r := IsValidIndex2(5, 0)
	assert.True(t, r)

	r = IsValidIndex2(5, 5)
	assert.False(t, r)

	r = IsValidIndex2(5, -1)
	assert.False(t, r)
}
