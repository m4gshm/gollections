package test

import (
	"testing"

	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/it"
	"github.com/stretchr/testify/assert"
)

func Test_Collect_Group(t *testing.T) {

	groups := collect.Group(it.ToKVIter[int, string](it.Of(K.V(1, "1"), K.V(2, "2"), K.V(2, "22"))))

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, []string{"1"}, groups[1])
	assert.Equal(t, []string{"2", "22"}, groups[2])
}

func Test_Collect_Map(t *testing.T) {

	groups := collect.Map(it.ToKVIter[int, string](it.Of(K.V(1, "1"), K.V(2, "2"), K.V(2, "22"))))

	assert.Equal(t, len(groups), 2)
	assert.Equal(t, "1", groups[1])
	assert.Equal(t, "22", groups[2])
}