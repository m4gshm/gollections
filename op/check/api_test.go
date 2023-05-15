package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValueTypeAlwaysNotNil(t *testing.T) {
	emptyString := ""
	ref := &emptyString
	assert.False(t, Nil(ref))
	ref = nil
	assert.True(t, Nil(ref))
}

func Test_NilIsNil(t *testing.T) {
	assert.False(t, NotNil[interface{}](nil))
	assert.True(t, Nil[interface{}](nil))
}

func Test_RealNil(t *testing.T) {
	assert.True(t, Nil[*string]((nil)))
	assert.True(t, Nil((*string)(nil)))
	assert.True(t, Nil((*int)(nil)))
	assert.True(t, Nil[*int](nil))
	assert.False(t, NotNil[*int](nil))
}

func Test_NilStruct(t *testing.T) {

	type someStruct struct{ somField string }

	v := someStruct{somField: "someValue"}

	r := &v

	assert.False(t, Nil(&v))
	assert.True(t, NotNil(r))
	r = nil
	assert.True(t, Nil(r))
	assert.True(t, Nil[*someStruct](nil))

	var i interface{}
	assert.False(t, Nil(&i))
}

func Test_ZeroStruct(t *testing.T) {
	type someStruct struct{ somField []string }

	var v someStruct
	assert.True(t, Zero(v))

}
