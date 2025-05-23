package seqexamples

import (
	"iter"
	"testing"

	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/stretchr/testify/assert"
)

func Test_SeqOf(t *testing.T) {

	var (
		ints  iter.Seq[int]          = seq.Of(1, 2, 3)
		pairs iter.Seq2[string, int] = seq2.OfMap(map[string]int{
			"first":  1,
			"second": 2,
			"third":  3,
		})
	)

	assert.Equal(t, []int{3, 2, 1}, sort.Desc(seq.Slice(seq2.Values(pairs))))
	assert.Equal(t, []string{"first", "second", "third"}, sort.Asc(seq.Slice(seq2.Keys(pairs))))
	assert.Equal(t, []int{1, 2, 3}, seq.Slice(ints))
}
