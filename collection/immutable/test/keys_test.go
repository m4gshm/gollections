package test

import (
	"testing"

	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/stretchr/testify/assert"
)

func Test_MapKeys_Zero_Safety(t *testing.T) {
	var collection immutable.MapKeys[int, string]

	collection.Head()
	collection.Convert(func(i int) int { return i })
	collection.Filter(func(_ int) bool { return true })
	collection.Slice()
	collection.Reduce(func(_, _ int) int { return 0 })
	s := collection.String()
	assert.Equal(t, "[]", s)

}

func Test_Map_Zero(t *testing.T) {
	var collection ordered.Map[int, string]

	collection.Loop()
	ptr.Of(collection.Head()).Next()
	collection.Convert(func(_ int, _ string) (int, string) { return 0, "" })
	collection.Filter(func(_ int, _ string) bool { return true })
	collection.Map()
	collection.Reduce(func(_, _ int, _, _ string) (int, string) { return 0, "" })
	s := collection.String()
	assert.Equal(t, "[]", s)

}
