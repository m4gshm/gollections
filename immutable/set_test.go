package immutable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	set := NewSet(1, 1, 2, 4, 3, 1)
	values := set.Values()

	assert.Equal(t, 4, set.Len())
	assert.Equal(t, 4, len(values))

}
