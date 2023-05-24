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
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	result, _ := breakkvloop.HasAny(kvl, func(key int, val string) bool { return key == 2 })

	assert.True(t, result)
}

func Test_HasAnyy(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	result, _ := breakkvloop.HasAnyy(kvl, func(key int, val string) (bool, error) { return key == 2, nil })

	assert.True(t, result)
}

func Test_Firstt(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	k, v, ok, _ := breakkvloop.Firstt(kvl, func(key int, val string) (bool, error) { return key == 2 || val == "three", nil })

	assert.True(t, ok)
	assert.Equal(t, 2, k)
	assert.Equal(t, "two", v)
}

func Test_Reduce(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	k, v, _ := breakkvloop.Reduce(kvl, func(kl, kr int, vl, vr string) (int, string) { return kl + kr, vl + vr })

	assert.Equal(t, 1+2+3, k)
	assert.Equal(t, "one"+"two"+"three", v)
}

func Test_Reducee(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "one"), k.V(2, "two"), k.V(3, "three"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	k, v, _ := breakkvloop.Reducee(kvl, func(kl, kr int, vl, vr string) (int, string, error) { return kl + kr, vl + vr, nil })

	assert.Equal(t, 1+2+3, k)
	assert.Equal(t, "one"+"two"+"three", v)
}

func Test_Convert(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	out, _ := breakkvloop.ToSlice(breakkvloop.Convert(kvl, func(k int, v string) (int, string) { return k * k, v + v }).Next, k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "11"), k.V(4, "22"), k.V(9, "33")), out)
}

func Test_Conv(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	out, _ := breakkvloop.ToSlice(breakkvloop.Conv(kvl, func(k int, v string) (int, string, error) { return k * k, v + v, nil }).Next, k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "11"), k.V(4, "22"), k.V(9, "33")), out)
}

func Test_Filter(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	out, _ := breakkvloop.ToSlice(breakkvloop.Filter(kvl, func(key int, val string) bool { return key != 2 }).Next, k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}

func Test_Filt(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	out, _ := breakkvloop.ToSlice(breakkvloop.Filt(kvl, func(key int, val string) (bool, error) { return key != 2, nil }).Next, k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}

func Test_Filt2(t *testing.T) {
	kvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	out, err := breakkvloop.ToSlice(breakkvloop.Filt(kvl, func(key int, val string) (bool, error) {
		ok := key <= 2
		return ok, op.IfElse(key == 2, errors.New("abort"), nil)
	}).Next, k.V[int, string])

	assert.Error(t, err)
	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(2, "2")), out)
}

func Test_To(t *testing.T) {
	bkvl := breakkvloop.From(loop.ToKV(slice.NewIter(slice.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))).Next, c.KV[int, string].Key, c.KV[int, string].Value).Next)

	kvl := breakkvloop.To(bkvl, func(err error) { assert.NoError(t, err) })

	out := kvloop.ToSlice(kvloop.Filter(kvl, func(key int, val string) bool { return key != 2 }).Next, k.V[int, string])

	assert.Equal(t, slice.Of(k.V(1, "1"), k.V(3, "3")), out)
}
