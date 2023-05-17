package boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/op/string_"
	"github.com/m4gshm/gollections/predicate/less"
	"github.com/m4gshm/gollections/predicate/match"
	"github.com/m4gshm/gollections/predicate/where"
	"github.com/m4gshm/gollections/slice"
)

func Test_Predicate(t *testing.T) {

	bob, _ := slice.First(users, where.Eq(User.Name, "Bob"))

	assert.Equal(t, "Bob", bob.Name())

}

func Test_ExtendedPredicate(t *testing.T) {

	userWithRoles := slice.Filter(users, match.To(User.Roles, slice.NotEmpty[[]Role]))

	assert.Equal(t, "Bob", userWithRoles[0].Name())
	assert.Equal(t, "Alice", userWithRoles[1].Name())

	userWithNamedRoles := slice.Filter(users, where.Any(User.Roles, match.To(Role.Name, string_.NotEmpty)))

	assert.Equal(t, "Bob", userWithNamedRoles[0].Name())
	assert.Equal(t, "Alice", userWithNamedRoles[1].Name())

	youngUsers := slice.Filter(users, where.Match(User.Age, less.OrEq(18)))

	assert.Equal(t, "Tom", youngUsers[0].Name())

}
