package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/match"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/predicate/one"
	"github.com/m4gshm/gollections/slice"
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

func Test_gtThan(t *testing.T) {
	gtC := more.Than("C")
	assert.False(t, gtC("A"))
	assert.False(t, gtC("B"))
	assert.False(t, gtC("C"))
	assert.True(t, gtC("D"))
	assert.True(t, gtC("E"))
}

func Test_gtOrEq(t *testing.T) {
	gtOrEqC := more.OrEq("C")
	assert.False(t, gtOrEqC("A"))
	assert.False(t, gtOrEqC("B"))
	assert.True(t, gtOrEqC("C"))
	assert.True(t, gtOrEqC("D"))
	assert.True(t, gtOrEqC("E"))
}

func Test_OneOf(t *testing.T) {
	oneOfACE := one.Of("A", "C", "E")
	assert.True(t, oneOfACE("A"))
	assert.False(t, oneOfACE("B"))
	assert.True(t, oneOfACE("C"))
	assert.False(t, oneOfACE("D"))
	assert.True(t, oneOfACE("E"))
}

func Test_Match_To(t *testing.T) {
	oneEvenMatcher := match.To(as.Is[int], func(i int) bool { return i%2 == 0 })

	assert.False(t, oneEvenMatcher(1))
	assert.True(t, oneEvenMatcher(2))

	anyEvenMatcher := match.To(as.Is, func(i []int) bool { return len(i) > 2 })

	assert.False(t, anyEvenMatcher(slice.Of(1, 2)))
	assert.True(t, anyEvenMatcher(slice.Of(1, 2, 3)))
}

func Test_Match_Any(t *testing.T) {
	anyEvenMatcher := match.Any(as.Is, func(i int) bool { return i%2 == 0 })

	assert.False(t, anyEvenMatcher(slice.Of(1, 3)))
	assert.True(t, anyEvenMatcher(slice.Of(1, 2)))
}
