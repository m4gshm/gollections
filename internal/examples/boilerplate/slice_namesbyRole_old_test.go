package boilerplate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamesByRole_Old(t *testing.T) {

	var namesByRole = map[string][]string{}
	add := func(role string, u User) {
		namesByRole[role] = append(namesByRole[role], u.Name())
	}
	for _, u := range users {
		roles := u.Roles()
		if len(roles) == 0 {
			add("", u)
		} else {
			for _, r := range roles {
				add(strings.ToLower(r.Name()), u)
			}
		}
	}

	assert.Equal(t, namesByRole[""], []string{"Tom"})
	assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
	assert.Equal(t, namesByRole["admin"], []string{"Bob"})

}
