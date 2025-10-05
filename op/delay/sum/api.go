// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice/sum"
)

// Of returns a sum builder function
func Of[T op.Summable](elements ...T) func() T {
	return func() T { return sum.Of(elements) }
}

// Over returns a sum builder function
func Over[T op.Summable](getters ...func() T) func() T {
	return func() T { return seq.Sum(seq.Convert(seq.Of(getters...), func(e func() T) T { return e() })) }
}
