package test

import (
	"testing"

	"github.com/m4gshm/gollections/predicate/break/predicate/less"
	"github.com/m4gshm/gollections/predicate/break/predicate/more"
	"github.com/m4gshm/gollections/predicate/break/predicate/one"
	"github.com/stretchr/testify/assert"
)

func Test_lessThan(t *testing.T) {
	lessC := less.Than("C")
	A, _ := lessC("A")
	assert.True(t, A)
	B, _ := lessC("B")
	assert.True(t, B)
	C, _ := lessC("C")
	assert.False(t, C)
	D, _ := lessC("D")
	assert.False(t, D)
}

func Test_lessOrEq(t *testing.T) {
	lessOrEqC := less.OrEq("C")
	A, _ := lessOrEqC("A")
	assert.True(t, A)
	B, _ := lessOrEqC("B")
	assert.True(t, B)
	C, _ := lessOrEqC("C")
	assert.True(t, C)
	D, _ := lessOrEqC("D")
	assert.False(t, D)
}

func Test_gtThan(t *testing.T) {
	gtC := more.Than("C")
	A, _ := gtC("A")
	assert.False(t, A)
	B, _ := gtC("B")
	assert.False(t, B)
	C, _ := gtC("C")
	assert.False(t, C)
	D, _ := gtC("D")
	assert.True(t, D)
	E, _ := gtC("E")
	assert.True(t, E)
}

func Test_gtOrEq(t *testing.T) {
	gtOrEqC := more.OrEq("C")
	A, _ := gtOrEqC("A")
	assert.False(t, A)
	B, _ := gtOrEqC("B")
	assert.False(t, B)
	C, _ := gtOrEqC("C")
	assert.True(t, C)
	D, _ := gtOrEqC("D")
	assert.True(t, D)
	E, _ := gtOrEqC("E")
	assert.True(t, E)
}

func Test_OneOf(t *testing.T) {
	oneOfACE := one.Of("A", "C", "E")
	A, _ := oneOfACE("A")
	assert.True(t, A)
	B, _ := oneOfACE("B")
	assert.False(t, B)
	C, _ := oneOfACE("C")
	assert.True(t, C)
	D, _ := oneOfACE("D")
	assert.False(t, D)
	E, _ := oneOfACE("E")
	assert.True(t, E)
}
