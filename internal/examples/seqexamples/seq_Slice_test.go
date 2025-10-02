package seqexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/seq"
)

func Test_ToSlice(t *testing.T) {

	filter := func(u User) bool { return u.age <= 30 }
	less30Names := seq.Convert(seq.Of(users...).Filter(filter), User.Name)

	names := seq.Slice(less30Names)
	//[Bob Tom]

	assert.Equal(t, []string{"Bob", "Tom"}, names)

	//or
	names = less30Names.Slice()
	//[Bob Tom]

	assert.Equal(t, []string{"Bob", "Tom"}, names)
}
