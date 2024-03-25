package loopexamples

import (
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/set"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/range_"
)

func Test_Iterating_Loop(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	next := uniques.Loop()
	for i, ok := next(); ok; i, ok = next() {
		doOp(i)
	}

}

func Test_Iterating_Iter(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	for iter, i, ok := uniques.First(); ok; i, ok = iter.Next() {
		doOp(i)
	}

}

func Test_Iterating_ForEach(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	uniques.ForEach(func(i int) { doOp(i) })

}

func Test_Iterating_For(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	uniques.For(func(i int) error {
		if i > 22 {
			return loop.Break
		}
		doOp(i)
		return nil
	})

}

func doOp(i int) {

}
