package loop

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/stretchr/testify/assert"
)

type next[T any] func() (element T, ok bool)

func Test_Slice_Vs_Loop(t *testing.T) {

	even := func(i int) bool { return i%2 == 0 }

	loopStream := loop.Convert(loop.Filter(loop.Of(1, 2, 3, 4), even).Next, strconv.Itoa)

	assert.Equal(t, []string{"2", "4"}, loop.Slice(loopStream.Next))

}
