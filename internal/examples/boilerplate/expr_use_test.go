package boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/expr/use"
	"github.com/m4gshm/gollections/op/delay/sum"
)

var isEmpty = func(s string) func() bool { return func() bool { return len(s) == 0 } }

func Test_UseSimple(t *testing.T) {

	user := User{name: "Bob", surname: "Smith"}

	fullName := use.If(len(user.surname) == 0, user.name).If(len(user.name) == 0, user.surname).
		ElseGet(func() string { return user.name + " " + user.surname })

	assert.Equal(t, "Bob Smith", fullName)

}

func Test_UseSimpleOld(t *testing.T) {

	user := User{name: "Bob", surname: "Smith"}

	fullName := ""
	if len(user.surname) == 0 {
		fullName = user.name
	} else if len(user.name) == 0 {
		fullName = user.surname
	} else {
		fullName = user.name + " " + user.surname
	}

	assert.Equal(t, "Bob Smith", fullName)

}

func Benchmark_UseIfElse(b *testing.B) {
	user := User{name: "Bob", surname: "Smith"}

	fullName := ""
	for i := 0; i < b.N; i++ {
		fullName = use.
			If(len(user.surname) == 0, user.name).
			If(len(user.name) == 0, user.surname).
			ElseGet(func() string { return user.name + " " + user.surname })
	}

	assert.Equal(b, "Bob Smith", fullName)
}

func Benchmark_UseIfElseGetSumOf(b *testing.B) {
	user := User{name: "Bob", surname: "Smith"}
	fullName := ""
	for i := 0; i < b.N; i++ {
		fullName = use.
			If(len(user.surname) == 0, user.name).
			If(len(user.name) == 0, user.surname).
			ElseGet(sum.Of(user.name, " ", user.surname))
	}

	assert.Equal(b, "Bob Smith", fullName)
}

func Benchmark_UseOtherElseGet(b *testing.B) {
	user := User{name: "Bob", surname: "Smith"}
	fullName := ""
	for i := 0; i < b.N; i++ {
		fullName = use.
			If(len(user.surname) == 0, user.name).
			Other(isEmpty(user.name), user.Surname).
			ElseGet(sum.Of(user.name, " ", user.surname))
	}

	assert.Equal(b, "Bob Smith", fullName)
}

func Benchmark_UseSimpleOld(b *testing.B) {
	user := User{name: "Bob", surname: "Smith"}
	fullName := ""
	for i := 0; i < b.N; i++ {
		if len(user.surname) == 0 {
			fullName = user.name
		} else if len(user.name) == 0 {
			fullName = user.surname
		} else {
			fullName = user.name + " " + user.surname
		}
	}
	assert.Equal(b, "Bob Smith", fullName)
}
