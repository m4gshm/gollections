package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/immutable/ordered/set"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
)

func Test_group_orderset_with_filtering_by_string_len(t *testing.T) {
	var groupedByLength = seq2.Group(seq.ToKV(set.Of(
		"seventh", "seventh",
		"first", "second", "third", "fourth",
		"fifth", "sixth", "eighth",
		"ninth", "tenth", "one", "two", "three", "1",
	).All, func(v string) int { return len(v) }, as.Is,
	).FilterKey(more.Than(3)).ConvertValue(func(v string) string { return v + "_" }))

	assert.Equal(t, []string{"first_", "third_", "fifth_", "sixth_", "ninth_", "tenth_", "three_"}, groupedByLength[5])
	assert.Equal(t, []string{"second_", "fourth_", "eighth_"}, groupedByLength[6])
	assert.Equal(t, []string{"seventh_"}, groupedByLength[7])
}
