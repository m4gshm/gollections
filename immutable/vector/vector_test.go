package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VectorIterate(t *testing.T) {
	v := Of(1, 2, 3, 4)
	i := 0
	for it := v.Iter(); it.HasNext(); {
		n := it.Next()
		expected, _ := v.Get(i)
		assert.Equal(t, expected, n)
		i++
	}
}
