package sliceexamples

import (
	"testing"

	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/stretchr/testify/assert"
)

// sorting copied slice API does not change the original slice

func Test_SortBy(t *testing.T) {

	var byName = sort.By(users, User.Name)
	//[{Alice 35 []} {Bob 26 []} {Chris 41 []} {Tom 18 []}]

	var byAgeReverse = sort.DescBy(users, User.Age)
	//[{Chris 41 []} {Alice 35 []} {Bob 26 []} {Tom 18 []}]

	assert.Equal(t, []User{
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Chris", age: 41},
		{name: "Tom", age: 18},
	}, byName)
	assert.Equal(t, []User{
		{name: "Chris", age: 41},
		{name: "Alice", age: 35},
		{name: "Bob", age: 26},
		{name: "Tom", age: 18},
	}, byAgeReverse)
}
