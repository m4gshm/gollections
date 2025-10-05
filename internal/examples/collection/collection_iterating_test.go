package collection

import (
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/set"
)

func Test_Iterating_ForEach(t *testing.T) {

	uniques := set.Of(1, 2, 3, 4, 5, 6)
	uniques.ForEach(doOp)

}

func doOp(i int) {

}
