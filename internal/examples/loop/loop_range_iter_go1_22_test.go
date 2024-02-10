package loop

import (
	"testing"

	"github.com/m4gshm/gollections/loop"
	"github.com/stretchr/testify/assert"
)

func Test_Range_Loop(t *testing.T) {

	var letters []rune
	for letter := range loop.RangeClosed('A', 'H').All {
		letters = append(letters, letter)
	}
	word := string(letters) //ABCDEFGH

	assert.Equal(t, "ABCDEFGH", word)

}
