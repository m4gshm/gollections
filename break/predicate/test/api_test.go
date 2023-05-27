package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/break/predicate"
	"github.com/m4gshm/gollections/break/predicate/eq"
	"github.com/m4gshm/gollections/break/predicate/less"
	"github.com/m4gshm/gollections/break/predicate/more"
)

func Test_Union(t *testing.T) {
	c, _ := predicate.Union[int]()(100)
	assert.False(t, c)
	c, _ = predicate.Union[int](predicate.Xor(eq.To(1), eq.To(1)))(1)
	assert.False(t, c)
	c, _ = predicate.Union[int](eq.To(1), less.Than(2))(1)
	assert.True(t, c)

	condition := predicate.Union[int](less.Than(3), more.Than(-1), predicate.Or[int](eq.To(0), eq.To(1)).Or(eq.To(2)))

	c, _ = condition(1)
	assert.True(t, c)
	c, _ = condition(0)
	assert.True(t, c)
	c, _ = condition(3)
	assert.False(t, c)
}
