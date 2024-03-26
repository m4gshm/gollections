//go:build goexperiment.rangefunc

package collection

import (
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/set"
	"github.com/m4gshm/gollections/loop/range_"
)

func Test_Iterating_Rangefunc(t *testing.T) {

	uniques := set.From(range_.Of(0, 100))
	for i := range uniques.All {
		doOp(i)
	}

}
