package mapexamples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

func Test_ToSlice(t *testing.T) {

	var employers = map[string]map[string]string{
		"devops": {"name": "Bob"},
		"jun":    {"name": "Tom"},
	}

	var users = map_.Slice(employers, func(title string, employer map[string]string) User {
		return User{name: employer["name"], roles: []Role{{name: title}}}
	})
	//[{name:Bob age:0 roles:[{name:devops}]} {name:Tom age:0 roles:[{name:jun}]}]

	fmt.Printf("%+v\n", users)

	assert.Equal(t,
		[]User{
			{name: "Bob", roles: []Role{{name: "devops"}}},
			{name: "Tom", roles: []Role{{name: "jun"}}},
		},
		slice.SortAsc(users, User.Name),
	)

}
