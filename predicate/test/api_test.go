package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
)

func Test_Xor(t *testing.T) {
	assert.False(t, predicate.Xor(eq.To(1), eq.To(1))(1))
	assert.False(t, predicate.Xor(eq.To(0), eq.To(0))(0))
	assert.True(t, predicate.Xor(eq.To(1), eq.To(0))(0))
	assert.True(t, predicate.Xor(eq.To(0), eq.To(1))(1))
}

func Test_Union(t *testing.T) {
	assert.False(t, predicate.Union[int]()(100))
	assert.False(t, predicate.Union(predicate.Xor(eq.To(1), eq.To(1)))(1))
	assert.True(t, predicate.Union(eq.To(1), less.Than(2))(1))

	condition := predicate.Union(less.Than(3), more.Than(-1), predicate.Or(eq.To(0), eq.To(1)).Or(eq.To(2)))

	assert.True(t, condition(1))
	assert.True(t, condition(0))
	assert.False(t, condition(3))
}
