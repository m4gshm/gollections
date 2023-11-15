package mapexamples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_/clone"
)

func Test_Clone(t *testing.T) {

	var bob = map[string]string{"name": "Bob"}
	var tom = map[string]string{"name": "Tom"}

	var employers = map[string]map[string]string{
		"devops": bob,
		"jun":    tom,
	}

	copy := clone.Of(employers)
	delete(copy, "jun")
	bob["name"] = "Superbob"

	fmt.Printf("%v\n", employers) //map[devops:map[name:Superbob] jun:map[name:Tom]]
	fmt.Printf("%v\n", copy)      //map[devops:map[name:Superbob]]

	assert.NotSame(t, copy, employers)

	assert.Equal(t, "Superbob", copy["devops"]["name"])

	assert.Contains(t, employers, "jun")
	assert.NotContains(t, copy, "jun")

}
