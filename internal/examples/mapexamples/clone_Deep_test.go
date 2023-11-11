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
		"":       tom,
	}

	copy := clone.Deep(employers, func(employer map[string]string) map[string]string {
		return clone.Of(employer)
	})
	delete(copy, "")
	copy["jun"] = tom
	bob["name"] = "Superbob"

	fmt.Printf("src  %v\n", employers) //src  map[:map[name:Tom] devops:map[name:Superbob]]
	fmt.Printf("copy %v\n", copy)      //copy map[devops:map[name:Bob] jun:map[name:Tom]]

	assert.NotSame(t, copy, employers)
	assert.Equal(t, "Bob", copy["devops"]["name"])

	assert.Contains(t, employers, "")
	assert.NotContains(t, employers, "jun")

	assert.NotContains(t, copy, "")
	assert.Contains(t, copy, "jun")

}
