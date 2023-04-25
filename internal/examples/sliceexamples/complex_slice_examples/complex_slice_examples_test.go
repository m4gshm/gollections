package sliceexamples

import (
	"strings"
	"testing"

	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/iter"
	sliceIter "github.com/m4gshm/gollections/iter/slice"
	"github.com/m4gshm/gollections/loop"
	loopConv "github.com/m4gshm/gollections/loop/convert"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/not"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/convert"
	"github.com/m4gshm/gollections/slice/group"
	sIter "github.com/m4gshm/gollections/slice/iter"
	"github.com/stretchr/testify/assert"
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
	usersByRole := group.InMultiple(users, func(u User) []string {
		return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
	})
	namesByRole := map_.ConvertValues(usersByRole, func(u []User) []string {
		return slice.Convert(u, User.Name)
	})

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
	// hasManager:=stream.Convert[User](Role.Name).HasAny(eq.To("Manager"))
	//new
	alice, ok := first.Of(users...).By(func(user User) bool {
		// return hasManager(user.Roles())
		return set.Convert(set.New(user.Roles()), Role.Name).HasAny(eq.To("Manager"))
	})

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

func Benchmark_FindFirsManager_Set(b *testing.B) {
	//new
	for i := 0; i < b.N; i++ {
		alice, ok := first.Of(users...).By(func(user User) bool {
			return set.Convert(set.New(user.Roles()), Role.Name).HasAny(eq.To("Manager"))
		})
		_, _ = alice, ok
	}
}

func Benchmark_FindFirsManager_Slice(b *testing.B) {
	//new
	for i := 0; i < b.N; i++ {
		alice, ok := first.Of(users...).By(func(user User) bool {
			return slice.Contains(slice.Convert(user.Roles(), Role.Name), "Manager")
		})
		_, _ = alice, ok
	}
}

func Benchmark_FindFirsManager_Loop_SliceIter(b *testing.B) {
	//new
	for i := 0; i < b.N; i++ {
		alice, ok := first.Of(users...).By(func(user User) bool {
			return loop.Contains(ptr.Of(sIter.Convert(user.Roles(), Role.Name)).Next, "Manager")
		})
		_, _ = alice, ok
	}
}

func Benchmark_FindFirsManager_Old(b *testing.B) {
	for i := 0; i < b.N; i++ {
		legacyAlice := User{}
		for _, u := range users {
			for _, r := range u.Roles() {
				if r.Name() == "Manager" {
					legacyAlice = u
				}
			}
		}
		_ = legacyAlice
	}
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

func Benchmark_AggregateFilteredRoles_Slice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		roles := slice.Flatt(users, User.Roles)
		roleNamesExceptManager := convert.AndFilter(roles, Role.Name, not.Eq("Manager"))
		_ = roleNamesExceptManager
	}
}

func Benchmark_AggregateFilteredRoles_Loop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		roles := sIter.Flatt(users, User.Roles)
		roleNamesExceptManager := loopConv.AndFilter(roles.Next, Role.Name, not.Eq("Manager"))
		_ = roleNamesExceptManager
	}
}

func Benchmark_AggregateFilteredRoles_Iter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		roles := sliceIter.Flatt(users, User.Roles)
		roleNamesExceptManager := iter.Filter(iter.Convert(roles, Role.Name), not.Eq("Manager"))
		_ = roleNamesExceptManager
	}
}

func Benchmark_AggregateFilteredRoles_Old(b *testing.B) {
	for i := 0; i < b.N; i++ {
		legacyRoleNamesExceptManager := []string{}
		for _, u := range users {
			for _, r := range u.Roles() {
				if n := r.Name(); n != "Manager" {
					legacyRoleNamesExceptManager = append(legacyRoleNamesExceptManager, n)
				}
			}
		}
		_ = legacyRoleNamesExceptManager
	}
}
