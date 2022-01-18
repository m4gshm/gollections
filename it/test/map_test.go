package it

import (
	"testing"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/stretchr/testify/assert"
)

func Test_KVIter_Iterate(t *testing.T) {
	type s struct {
		name string
	}

	values := map[string]*s{"first": {"first_name"}, "second": {"first_second"}}
	result := map[string]*s{}

	for it := it.NewKV(values); it.HasNext(); {
		k, v, _ := it.Get()
		result[k] = v
	}
	assert.Equal(t, len(values), len(values))

	for k, v := range values {
		assert.Equal(t, v, result[k])
	}
}
