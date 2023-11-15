package mapexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_"
)

func Test_ConvertValues(t *testing.T) {

	var employers = map[string]map[string]string{
		"devops": {"name": "Bob"},
		"jun":    {"name": "Tom"},
	}

	var vals = map_.ConvertValues(employers, func(employer map[string]string) string {
		return employer["name"]
	})
	//map[devops:Bob jun:Tom]

	assert.Equal(t, map[string]string{"devops": "Bob", "jun": "Tom"}, vals)
}
