package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertByAndChain(t *testing.T) {

	var toString Converter[int, string] = func(i int) (string, error) { return fmt.Sprintf("%d", i), nil }
	var addTail Converter[string, string] = func(s string) (string, error) { return s + "_tail", nil }

	converted, err := Convert(Of(1, 2, 3), And(toString, addTail))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Of("1_tail", "2_tail", "3_tail"), converted)
}
