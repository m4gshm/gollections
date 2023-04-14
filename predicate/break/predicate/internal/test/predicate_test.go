package test

import (
	"testing"

	"github.com/m4gshm/gollections/predicate/break/predicate"
	"github.com/stretchr/testify/assert"
)

func Test_XorFalse(t *testing.T) {
	xor1 := predicate.Xor(predicate.Eq("A"), predicate.Eq("A"))
	ok1, _ := xor1("A")
	assert.False(t, ok1)

	xor2 := predicate.Xor(predicate.Not(predicate.Eq("A")), predicate.Not(predicate.Eq("A")))
	ok2, _ := xor2("A")
	assert.False(t, ok2)

	xor3 := predicate.Xor(predicate.Eq("A"), predicate.Eq("B"))
	ok3, _ := xor3("A")
	assert.True(t, ok3)

	xor4 := predicate.Xor(predicate.Eq("B"), predicate.Eq("A"))
	ok4, _ := xor4("A")
	assert.True(t, ok4)
}
