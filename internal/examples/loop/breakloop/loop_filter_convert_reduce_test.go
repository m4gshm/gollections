package breakableloop

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/break/loop"
	"github.com/stretchr/testify/assert"
)

type (
	next[T any]      func() (element T, ok bool, err error)
	kvNext[K, V any] func() (key K, value V, ok bool, err error)
)

func Test_Slice_Vs_Loop(t *testing.T) {

	iter := loop.Conv(loop.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
	result, err := loop.Slice(iter.Next)

	assert.Equal(t, []int{1, 2, 3}, result)
	assert.ErrorContains(t, err, "invalid syntax")

}
