//go:build goexperiment.rangefunc

package collection

import (
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/set"
)

func Test_Iterating_Rangefunc(t *testing.T) {

	uniques := set.Of(1, 2, 3, 4, 5, 6)
	for i := range uniques.All {
		doOp(i)
	}

}
