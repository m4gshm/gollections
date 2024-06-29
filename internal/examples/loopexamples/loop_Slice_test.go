package loopexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/loop"
)

func Test_ToSlice(t *testing.T) {

	filter := func(u User) bool { return u.age <= 30 }
	names := loop.Slice(loop.Convert(loop.Filter(loop.Of(users...), filter), User.Name))
	//[Bob Tom]

	assert.Equal(t, []string{"Bob", "Tom"}, names)
}
