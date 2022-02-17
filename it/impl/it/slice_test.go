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
