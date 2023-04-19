package test

import (
	"testing"

	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/stretchr/testify/assert"
)

func Test_lessThan(t *testing.T) {
	lessC := less.Than("C")
	assert.True(t, lessC("A"))
	assert.True(t, lessC("B"))
	assert.False(t, lessC("C"))
	assert.False(t, lessC("D"))
}

func Test_lessOrEq(t *testing.T) {
	lessOrEqC := less.OrEq("C")
	assert.True(t, lessOrEqC("A"))
	assert.True(t, lessOrEqC("B"))
	assert.True(t, lessOrEqC("C"))
	assert.False(t, lessOrEqC("D"))
}

func Test_moreThan(t *testing.T) {
	moreC := more.Than("C")
	assert.False(t, moreC("A"))
	assert.False(t, moreC("B"))
	assert.False(t, moreC("C"))
	assert.True(t, moreC("D"))
	assert.True(t, moreC("E"))
}

func Test_moreOrEq(t *testing.T) {
	moreOrEqC := more.OrEq("C")
	assert.False(t, moreOrEqC("A"))
	assert.False(t, moreOrEqC("B"))
	assert.True(t, moreOrEqC("C"))
	assert.True(t, moreOrEqC("D"))
	assert.True(t, moreOrEqC("E"))
}

func Test_OneOf(t *testing.T) {
	oneOfACE := one.Of("A", "C", "E")
	assert.True(t, oneOfACE("A"))
	assert.False(t, oneOfACE("B"))
	assert.True(t, oneOfACE("C"))
	assert.False(t, oneOfACE("D"))
	assert.True(t, oneOfACE("E"))
}
