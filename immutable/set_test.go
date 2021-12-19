package immutable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	values := NewSet(1, 1, 2, 4, 3, 1).Values()

	assert.Equal(t, 4, NewSet(1, 1, 2, 4, 3, 1).Len())
	assert.Equal(t, 4, len(values))

}
