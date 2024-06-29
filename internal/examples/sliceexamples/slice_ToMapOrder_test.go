package sliceexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
)

func Test_Slice_ToMapOrder(t *testing.T) {

	var names, agePerName = slice.MapOrder(users, User.Name, User.Age)

	//"[Bob Alice Tom Chris]"
	//"map[Alice:35 Bob:26 Chris:41 Tom:18]"

	assert.Equal(t, []string{"Bob", "Alice", "Tom", "Chris"}, names)
	assert.NotNil(t, agePerName)
}
