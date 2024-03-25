package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	breakkvloop "github.com/m4gshm/gollections/break/kv/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/k"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func Test_HasAny(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value)

	result := kvloop.HasAny(kvl, func(key int, _ string) bool { return key == 2 })

	assert.True(t, result)
}

func Test_Firstt(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value)

	k, v, ok, _ := kvloop.Firstt(kvl, func(key int, val string) (bool, error) { return key == 2 || val == "three", nil })

	assert.True(t, ok)
	assert.Equal(t, 2, k)
	assert.Equal(t, "two", v)
}

func Test_Reduce(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value)

	k, v := kvloop.Reduce(kvl, func(kl, kr int, vl, vr string) (int, string) { return kl + kr, vl + vr })

	assert.Equal(t, 1+2+3, k)
	assert.Equal(t, "one"+"two"+"three", v)
}

func Test_Reducee(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value)

	k, v, _ := kvloop.Reducee(kvl, func(kl, kr int, vl, vr string) (int, string, error) { return kl + kr, vl + vr, nil })

	assert.Equal(t, 1+2+3, k)
	assert.Equal(t, "one"+"two"+"three", v)
}

func Test_Convert(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value)

	out := kvloop.ToSlice(kvloop.Convert(kvl, func(k int, v string) (int, string) { return k * k, v + v }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "11"), k.V(4, "22"), k.V(9, "33")), out)
}

func Test_Conv(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value)

	out, _ := breakkvloop.ToSlice(kvloop.Conv(kvl, func(k int, v string) (int, string, error) { return k * k, v + v, nil }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "11"), k.V(4, "22"), k.V(9, "33")), out)
}

func Test_Filter(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value)

	out := kvloop.ToSlice(kvloop.Filter(kvl, func(key int, _ string) bool { return key != 2 }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}

func Test_Filt(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value)

	out, _ := breakkvloop.ToSlice(kvloop.Filt(kvl, func(key int, _ string) (bool, error) { return key != 2, nil }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}

func Test_Filt2(t *testing.T) {
	kvl := loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value)

	out, err := breakkvloop.ToSlice(kvloop.Filt(kvl, func(key int, _ string) (bool, error) {
		ok := key <= 2
		return ok, op.IfElse(key == 2, errors.New("abort"), nil)
	}), k.V[int, string])

	assert.Error(t, err)
	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(2, "2")), out)
}

func Test_NewIter(t *testing.T) {
	s := slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
	i := 0
	loop := kvloop.New(s, func(s []c.KV[int, string]) bool { return i < len(s) }, func(s []c.KV[int, string]) (int, string) {
		n := s[i]
		i++
		return n.K, n.V
	})

	out := kvloop.ToSlice(loop, k.V[int, string])
	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), out)
}
