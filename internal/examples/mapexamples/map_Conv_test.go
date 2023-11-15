package mapexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_"
)

func Test_Conv(t *testing.T) {

	var employers = map[string]map[string]string{
		"devops": {"name": "Bob"},
		"jun":    {"name": "Tom"},
	}

	var all, err = map_.Conv(employers, func(title string, employer map[string]string) (string, string, error) {
		return string([]rune(title)[0]), employer["name"], nil
	})
	//map[d:Bob j:Tom], nil

	assert.Equal(t, map[string]string{"d": "Bob", "j": "Tom"}, all)
	assert.Nil(t, err)

}
