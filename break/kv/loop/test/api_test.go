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
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value))

	result, _ := breakkvloop.HasAny(kvl, func(key int, _ string) bool { return key == 2 })

	assert.True(t, result)
}

func Test_HasAnyy(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value))

	result, _ := breakkvloop.HasAnyy(kvl, func(key int, _ string) (bool, error) { return key == 2, nil })

	assert.True(t, result)
}

func Test_Firstt(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value))

	k, v, ok, _ := breakkvloop.Firstt(kvl, func(key int, val string) (bool, error) { return key == 2 || val == "three", nil })

	assert.True(t, ok)
	assert.Equal(t, 2, k)
	assert.Equal(t, "two", v)
}

func Test_Reduce(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value))

	k, v, ok, _ := breakkvloop.ReduceOK(kvl, func(kl, kr int, vl, vr string) (int, string) { return kl + kr, vl + vr })

	assert.True(t, ok)
	assert.Equal(t, 1+2+3, k)
	assert.Equal(t, "one"+"two"+"three", v)
}

func Test_Reduce_Empty(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of[c.KV[int, string]](), c.KV[int, string].Key, c.KV[int, string].Value))

	_, _, ok, _ := breakkvloop.ReduceOK(kvl, func(kl, kr int, vl, vr string) (int, string) { return kl + kr, vl + vr })

	assert.False(t, ok)
}

func Test_Reduce_Nil(t *testing.T) {
	var l loop.Loop[c.KV[int, string]]
	kvl := breakkvloop.From(loop.KeyValue(l, c.KV[int, string].Key, c.KV[int, string].Value))

	_, _, ok, _ := breakkvloop.ReduceOK(kvl, func(kl, kr int, vl, vr string) (int, string) { return kl + kr, vl + vr })

	assert.False(t, ok)
}

func Test_Reducee(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three")), c.KV[int, string].Key, c.KV[int, string].Value))

	k, v, ok, _ := breakkvloop.ReduceeOK(kvl, func(kl, kr int, vl, vr string) (int, string, error) { return kl + kr, vl + vr, nil })

	assert.True(t, ok)
	assert.Equal(t, 1+2+3, k)
	assert.Equal(t, "one"+"two"+"three", v)
}

func Test_Reducee_Empty(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of[c.KV[int, string]](), c.KV[int, string].Key, c.KV[int, string].Value))

	_, _, ok, _ := breakkvloop.ReduceeOK(kvl, func(kl, kr int, vl, vr string) (int, string, error) { return kl + kr, vl + vr, nil })

	assert.False(t, ok)
}

func Test_Reducee_Nil(t *testing.T) {
	var l loop.Loop[c.KV[int, string]]
	kvl := breakkvloop.From(loop.KeyValue(l, c.KV[int, string].Key, c.KV[int, string].Value))

	_, _, ok, _ := breakkvloop.ReduceeOK(kvl, func(kl, kr int, vl, vr string) (int, string, error) { return kl + kr, vl + vr, nil })

	assert.False(t, ok)
}

func Test_Convert(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value))

	out, _ := breakkvloop.Slice(breakkvloop.Convert(kvl, func(k int, v string) (int, string) { return k * k, v + v }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "11"), k.V(4, "22"), k.V(9, "33")), out)
}

func Test_Conv(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value))

	out, _ := breakkvloop.Slice(breakkvloop.Conv(kvl, func(k int, v string) (int, string, error) { return k * k, v + v, nil }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "11"), k.V(4, "22"), k.V(9, "33")), out)
}

func Test_Filter(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value))

	out, _ := breakkvloop.Slice(breakkvloop.Filter(kvl, func(key int, _ string) bool { return key != 2 }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}

func Test_Filt(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value))

	out, _ := breakkvloop.Slice(breakkvloop.Filt(kvl, func(key int, _ string) (bool, error) { return key != 2, nil }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}

func Test_Filt2(t *testing.T) {
	kvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value))

	out, err := breakkvloop.Slice(breakkvloop.Filt(kvl, func(key int, _ string) (bool, error) {
		ok := key <= 2
		return ok, op.IfElse(key == 2, errors.New("abort"), nil)
	}), k.V[int, string])

	assert.Error(t, err)
	assert.Equal(t, slice.Of(k.V(1, "1")), out)
}

func Test_To(t *testing.T) {
	bkvl := breakkvloop.From(loop.KeyValue(loop.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3")), c.KV[int, string].Key, c.KV[int, string].Value))

	kvl := breakkvloop.To(bkvl, func(err error) { assert.NoError(t, err) })

	out := kvloop.Slice(kvloop.Filter(kvl, func(key int, _ string) bool { return key != 2 }), k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}
