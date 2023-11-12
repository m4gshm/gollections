package mapexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_"
)

func Test_ConvertKeys(t *testing.T) {

	var employers = map[string]map[string]string{
		"devops": {"name": "Bob"},
		"jun":    {"name": "Tom"},
	}

	var keys = map_.ConvertKeys(employers, func(title string) string {
		return string([]rune(title)[0])
	})
	//map[d:map[name:Bob] j:map[name:Tom]]

	assert.Equal(t, map[string]map[string]string{"d": {"name": "Bob"}, "j": {"name": "Tom"}}, keys)

}
