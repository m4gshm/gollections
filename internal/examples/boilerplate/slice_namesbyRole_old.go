package boilerplate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _(t *testing.T) {

	namesByRole := map[string][]string{}
	for _, u := range users {
		roles := u.Roles()
		if len(roles) == 0 {
			lr := ""
			names := namesByRole[lr]
			names = append(names, u.Name())
			namesByRole[lr] = names
		} else {
			for _, r := range roles {
				lr := strings.ToLower(r.Name())
				names := namesByRole[lr]
				names = append(names, u.Name())
				namesByRole[lr] = names
			}
		}
	}

	assert.Equal(t, namesByRole[""], []string{"Tom"})
	assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
	assert.Equal(t, namesByRole["admin"], []string{"Bob"})
	
}
