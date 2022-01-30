package test

import (
	"fmt"
	"testing"

	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/stretchr/testify/assert"
)

func Test_StringRepresentation(t *testing.T) {
	expected := fmt.Sprint(slice.Of(1, 2, 3, 4))

	var i1, i2, i3, i4 = 1, 2, 3, 4
	order := slice.Of(&i4, &i3, &i2, &i1)

	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}

	expected = "[4:4 3:3 2:2 1:1]"
	actual := map_.ToStringOrdered(order, elements)

	assert.Equal(t, expected, actual)
}
