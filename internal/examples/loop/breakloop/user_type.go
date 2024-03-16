package breakableloop

type User struct {
	name, surname string
	age           int
	roles         []Role
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

func (u User) Surname() string {
	return u.surname
}

func (u User) Age() int {
	return u.age
}

func (u User) Roles() []Role {
	return u.roles
}
