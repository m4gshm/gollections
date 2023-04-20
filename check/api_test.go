package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValueTypeAlwaysNotNil(t *testing.T) {
	emptyString := ""
	ref := &emptyString
	assert.Equal(t, false, Nil(ref))
	ref = nil
	assert.Equal(t, true, Nil(ref))
}

func Test_NilIsNil(t *testing.T) {
	assert.Equal(t, false, NotNil[interface{}](nil))
	assert.Equal(t, true, Nil[interface{}](nil))
}

func Test_RealNil(t *testing.T) {
	assert.Equal(t, true, Nil[*string]((nil)))
	assert.Equal(t, true, Nil((*string)(nil)))
	assert.Equal(t, true, Nil((*int)(nil)))
	assert.Equal(t, true, Nil[*int](nil))
	assert.Equal(t, false, NotNil[*int](nil))
}

func Test_NilStruct(t *testing.T) {

	type someStruct struct{ somField string }

	v := someStruct{somField: "someValue"}

	r := &v

	assert.Equal(t, false, Nil(&v))
	assert.Equal(t, true, NotNil(r))
	r = nil
	assert.Equal(t, true, Nil(r))
	assert.Equal(t, true, Nil[*someStruct](nil))

	var i interface{}
	assert.Equal(t, false, Nil(&i))
}

func Test_Empty(t *testing.T) {
	s := ""
	assert.Equal(t, true, Empty(s))

	i := []int{}
	assert.Equal(t, true, Empty(i))
}
