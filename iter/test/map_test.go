package it

import (
	"testing"

	"github.com/m4gshm/gollections/map_/iter"
	"github.com/stretchr/testify/assert"
)

func Test_NewEmbedMapKV_Iterate(t *testing.T) {
	type s struct {
		name string
	}

	values := map[string]*s{"first": {"first_name"}, "second": {"first_second"}}
	result := map[string]*s{}
	it := iter.New(values)
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		result[k] = v
	}
	assert.Equal(t, len(values), len(values))

	for k, v := range values {
		assert.Equal(t, v, result[k])
	}
}
