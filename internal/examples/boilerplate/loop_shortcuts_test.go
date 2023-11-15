package boilerplate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/group"
	"github.com/m4gshm/gollections/slice"
)

func Test_Loop_schortcuts(t *testing.T) {

	data := loop.Of("Bob", "Chris", "Alice")

	lowers := loop.Convert(data, strings.ToLower)

	var lengthMap map[int][]string = group.Of(lowers.Next, func(name string) int { return len(name) }, as.Is[string]) // converting to map

	assert.Equal(t, slice.Of("chris", "alice"), lengthMap[5])

}
