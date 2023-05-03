package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/collection/group"
	"github.com/m4gshm/gollections/collection/immutable/oset"
	"github.com/m4gshm/gollections/predicate/more"
)

func Test_group_orderset_with_filtering_by_string_len(t *testing.T) {

	var groupedByLength = group.Of(oset.Of(
		"seventh", "seventh", //duplicate
		"first", "second", "third", "fourth",
		"fifth", "sixth", "eighth",
		"ninth", "tenth", "one", "two", "three", "1",
		"second", //duplicate
	), func(v string) int { return len(v) },
	).FilterKey(
		more.Than(3),
	).ConvertValue(
		func(v string) string { return v + "_" },
	).Map()

	assert.Equal(t, map[int][]string{
		5: {"first_", "third_", "fifth_", "sixth_", "ninth_", "tenth_", "three_"},
		6: {"second_", "fourth_", "eighth_"},
		7: {"seventh_"},
	}, groupedByLength)

}
