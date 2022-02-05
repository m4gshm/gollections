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

func Test_Vector_Sort(t *testing.T) {
	var (
		elements = Of(3, 1, 5, 6, 8, 0, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, Of(-2, 0, 1, 3, 5, 6, 8), sorted)
}
