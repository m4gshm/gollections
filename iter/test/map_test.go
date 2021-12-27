package iter

import (
	"testing"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/stretchr/testify/assert"
)

func Test_MapIterator(t *testing.T) {
	type s struct {
		name string
	}

	values := map[string]*s{"first": {"first_name"}, "second": {"first_second"}}
	result := map[string]*s{}

	for it := iter.NewKV(values); it.HasNext(); {
		k, v := it.Get().Get()
		result[k] = v
	}
	assert.Equal(t, len(values), len(values))

	for k, v := range values {
		assert.Equal(t, v, result[k])
	}
}

func Test_MapIterateResetIterat(t *testing.T) {
	type s struct {
		name string
	}

	values := map[string]*s{"first": {"first_name"}, "second": {"first_second"}}
	result1 := map[string]*s{}

	it := iter.NewKV(values)
	for it.HasNext() {
		k, v := it.Get().Get()
		result1[k] = v
	}
	it.Reset()

	result2 := map[string]*s{}
	for it.HasNext() {
		k, v := it.Get().Get()
		result2[k] = v
	}

	for k, v := range values {
		assert.Equal(t, v, result1[k])
		assert.Equal(t, v, result2[k])
	}
}
