//go:build goexperiment.rangefunc

package loopexamples

import (
	"testing"

	"github.com/m4gshm/gollections/loop/range_"
)

func Test_Iterating_Rangefunc(t *testing.T) {

	for i := range range_.Of(0, 100).All {
		doOp(i)
	}

}
