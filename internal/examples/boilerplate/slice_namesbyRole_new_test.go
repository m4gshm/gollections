package boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import (
	"strings"
	"github.com/m4gshm/gollections/slice/convert"
	"github.com/m4gshm/gollections/slice/group"
)

func Test_NamesByRole_New(t *testing.T) {

	var namesByRole = group.ByMultipleKeys(users, func(u User) []string {
		return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
	}, User.Name)
	// map[:[Tom] admin:[Bob] manager:[Bob Alice]]

	assert.Equal(t, namesByRole[""], []string{"Tom"})
	assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
	assert.Equal(t, namesByRole["admin"], []string{"Bob"})

}
