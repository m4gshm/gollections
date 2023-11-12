package mapexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_"
)

func Test_Convert(t *testing.T) {

	var employers = map[string]map[string]string{
		"devops": {"name": "Bob"},
		"jun":    {"name": "Tom"},
	}

	var all = map_.Convert(employers, func(title string, employer map[string]string) (string, string) {
		return string([]rune(title)[0]), employer["name"]
	})
	//map[d:Bob j:Tom]

	assert.Equal(t, map[string]string{"d": "Bob", "j": "Tom"}, all)

}
