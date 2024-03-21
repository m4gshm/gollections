package breakableloop

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/break/loop"
)

func Test_Slice_Vs_Loop(t *testing.T) {

	intSeq := loop.Conv(loop.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
	ints, err := loop.Slice(intSeq)

	assert.Equal(t, []int{1, 2, 3}, ints)
	assert.ErrorContains(t, err, "invalid syntax")

}
