package mapexamples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_/clone"
)

func Test_CloneDeep(t *testing.T) {

	var bob = map[string]string{"name": "Bob"}
	var tom = map[string]string{"name": "Tom"}

	var employers = map[string]map[string]string{
		"devops": bob,
		"jun":    tom,
	}

	copy := clone.Deep(employers, func(employer map[string]string) map[string]string {
		return clone.Of(employer)
	})
	delete(copy, "jun")
	bob["name"] = "Superbob"

	fmt.Printf("%v\n", employers) //map[devops:map[name:Superbob] jun:map[name:Tom]]
	fmt.Printf("%v\n", copy)      //map[devops:map[name:Bob]]

	assert.NotSame(t, &copy, &employers)

	assert.Equal(t, "Bob", copy["devops"]["name"])

	assert.Contains(t, employers, "jun")
	assert.NotContains(t, copy, "jun")

}
