package mapexamples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
)

func Test_Keys(t *testing.T) {

	var employers = map[string]map[string]string{
		"devops": {"name": "Bob"},
		"jun":    {"name": "Tom"},
	}

	keys := map_.Keys(employers)     //[devops jun]
	values := map_.Values(employers) //[map[name:Bob] map[name:Tom]]

	assert.Equal(t, slice.Of("devops", "jun"), sort.Asc(keys))

	names := slice.Flat(slice.Convert(values, map_.Values), as.Is[[]string])
	assert.Equal(t, slice.Of("Bob", "Tom"), sort.Asc(names))

}
