package iter

import (
	"testing"

	"github.com/m4gshm/container/iter"
	"github.com/stretchr/testify/assert"
)

func Test_MapIterator(t *testing.T) {
	type s struct {
		name string
	}

	values := map[string]*s{"first": {"first_name"}, "second": {"first_second"}}
	result := map[string]*s{}
	it := iter.WrapMap(values)
	for it.HasNext() {
		kv := it.Get()
		result[kv.Key()] = kv.Value()
	}
	assert.Equal(t, len(values), len(values))

	for k, v := range values {
		assert.Equal(t, v, result[k])
	}
}
