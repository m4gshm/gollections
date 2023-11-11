package sliceexamples

type User struct {
	name  string
	age   int
	roles []Role
}

type Role struct {
	name string
}

func (u Role) Name() string {
	return u.name
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