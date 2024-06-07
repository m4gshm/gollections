package sliceexamples

import (
	"testing"

	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_Slice_ToMap(t *testing.T) {

	var agePerGroup = slice.ToMap(users, User.Name, User.Age)

	//"map[Alice:35 Bob:26 Chris:41 Tom:18]"

	assert.Equal(t, 35, agePerGroup["Alice"])
	assert.Equal(t, 18, agePerGroup["Tom"])
}
