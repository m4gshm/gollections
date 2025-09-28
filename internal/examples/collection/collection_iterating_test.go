package collection

import (
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/set"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/range_"
)

func Test_Iterating_ForEach(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	uniques.ForEach(doOp)

}

func Test_Iterating_For(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	uniques.For(func(i int) error {
		if i > 22 {
			return loop.Break
		}
		doOp(i)
		return loop.Continue
	})

}

func doOp(i int) {

}
