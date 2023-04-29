package test

import (
	"testing"

	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/ptr"
	"github.com/stretchr/testify/assert"
)

func Test_MapKeys_Zero_Safety(t *testing.T) {
	var collection immutable.MapKeys[int, string]

	collection.Iter()
	collection.Head()
	collection.Convert(func(i int) int { return i })
	collection.Filter(func(i int) bool { return true })
	collection.Slice()
	collection.Reduce(func(i1, i2 int) int { return 0 })
	s := collection.String()
	assert.Equal(t, "[]", s)

}

func Test_Map_Zero(t *testing.T) {
	var collection ordered.Map[int, string]

	collection.Iter().Next()
	ptr.Of(collection.Head()).Next()
	collection.Convert(func(i int, s string) (int, string) { return 0, "" })
	collection.Filter(func(i int, s string) bool { return true })
	collection.Map()
	collection.Reduce(func(i1 int, s1 string, i2 int, s2 string) (int, string) { return 0, "" })
	s := collection.String()
	assert.Equal(t, "[]", s)

}
