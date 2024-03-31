package loopexamples

import (
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/range_"
)

func Test_Iterating_EmbeddedFor(t *testing.T) {

	next := range_.Of(0, 100)
	for i, ok := next(); ok; i, ok = next() {
		doOp(i)
	}

}

func Test_Iterating_EmbeddedFor2(t *testing.T) {

	for next, i, ok := range_.Of(0, 100).Crank(); ok; i, ok = next() {
		doOp(i)
	}

}

func Test_Iterating_ForEach(t *testing.T) {

	range_.Of(0, 100).ForEach(doOp)

}

func Test_Iterating_For(t *testing.T) {

	range_.Of(0, 100).For(func(i int) error {
		if i > 22 {
			return loop.Break
		}
		doOp(i)
		return loop.Continue
	})

}

func doOp(i int) {

}
