package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValueTypeAlwaysNotNil(t *testing.T) {
	assert.Equal(t, true, NotNil("string"))
	assert.Equal(t, false, Nil("string"))
	assert.Equal(t, true, NotNil(1.0))
	assert.Equal(t, false, Nil(1.0))
}

func Test_NilIsNil(t *testing.T) {
	assert.Equal(t, false, NotNil[interface{}](nil))
	assert.Equal(t, true, Nil[interface{}](nil))
}

func Test_InterfcaeAsNilOrNotNil(t *testing.T) {
	var i interface{}
	assert.Equal(t, false, NotNil(i))
	assert.Equal(t, true, Nil(i))

	i = make([]int, 0)

	assert.Equal(t, true, NotNil(i))
	assert.Equal(t, false, Nil(i))
}

func Test_MapAsNilOrNotNil(t *testing.T) {
	var m map[string]string
	assert.Equal(t, true, Nil(m))
	m2 := map[string]string{}
	assert.Equal(t, true, NotNil(m2))

}
