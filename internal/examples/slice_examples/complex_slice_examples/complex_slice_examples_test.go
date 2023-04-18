package slice_examples

import (
	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/predicate/not"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/convert"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type User struct {
	name  string
	age   int
	roles []Role
}

type Role struct {
	name string
}

func (r Role) Name() string {
	return r.name
}

func (u User) Name() string {
	return u.name
}

func (u User) Age() int {
	return u.age
}

func (u User) Roles() []Role {
	return u.roles
}

var users = []User{
	{name: "Bob", age: 26, roles: []Role{{"Admin"}, {"manager"}}},
	{name: "Alice", age: 35, roles: []Role{{"Manager"}}},
	{name: "Tom", age: 18},
}

func Test_GroupBySeveralKeysAndConvertMapValues(t *testing.T) {

	//new
	usersByRole := group.InMultiple(users, func(u User) []string { return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower) })
	namesByRole := map_.ConvertValues(usersByRole, func(u []User) []string { return slice.Convert(u, User.Name) })

	assert.Equal(t, namesByRole[""], []string{"Tom"})
	assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
	assert.Equal(t, namesByRole["admin"], []string{"Bob"})

	//old
	legacyNamesByRole := map[string][]string{}
	for _, u := range users {
		roles := u.Roles()
		if len(roles) == 0 {
			lr := ""
			names := legacyNamesByRole[lr]
			names = append(names, u.Name())
			legacyNamesByRole[lr] = names
		} else {
			for _, r := range roles {
				lr := strings.ToLower(r.Name())
				names := legacyNamesByRole[lr]
				names = append(names, u.Name())
				legacyNamesByRole[lr] = names
			}
		}
	}

	assert.Equal(t, legacyNamesByRole[""], []string{"Tom"})
	assert.Equal(t, legacyNamesByRole["manager"], []string{"Bob", "Alice"})
	assert.Equal(t, legacyNamesByRole["admin"], []string{"Bob"})
}

func Test_FindFirsManager(t *testing.T) {
	//new
	alice, ok := first.Of(users...).By(func(u User) bool { return set.New(slice.Convert(u.Roles(), Role.Name)).Contains("Manager") })

	assert.True(t, ok)
	assert.Equal(t, "Alice", alice.Name())

	//plain old
	legacyAlice := User{}
	for _, u := range users {
		for _, r := range u.Roles() {
			if r.Name() == "Manager" {
				legacyAlice = u

			}
		}
	}
	assert.Equal(t, "Alice", legacyAlice.Name())
}

func Test_AggregateFilteredRoles(t *testing.T) {
	//new
	roles := slice.Flatt(users, User.Roles)
	roleNamesExceptManager := convert.AndFilter(roles, Role.Name, not.Eq("Manager"))

	assert.Equal(t, slice.Of("Admin", "manager"), roleNamesExceptManager)

	//plain old
	legacyRoleNamesExceptManager := []string{}
	for _, u := range users {
		for _, r := range u.Roles() {
			if n := r.Name(); n != "Manager" {
				legacyRoleNamesExceptManager = append(legacyRoleNamesExceptManager, n)
			}
		}
	}
	assert.Equal(t, slice.Of("Admin", "manager"), legacyRoleNamesExceptManager)
}