package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_ToSlice(t *testing.T) {

	filter := func(u User) bool { return u.age <= 30 }
	names := seq.Slice(seq.Convert(seq.Filter(seq.Of(users...), filter), User.Name))
	//[Bob Tom]

	assert.Equal(t, []string{"Bob", "Tom"}, names)
}
