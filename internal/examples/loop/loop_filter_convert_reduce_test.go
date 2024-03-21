package loop

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_Slice_Vs_Loop(t *testing.T) {

	even := func(i int) bool { return i%2 == 0 }
	stringSeq := loop.Convert(loop.Filter(loop.Of(1, 2, 3, 4), even), strconv.Itoa)

	assert.Equal(t, []string{"2", "4"}, loop.Slice(stringSeq))

}
