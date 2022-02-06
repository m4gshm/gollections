package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VectorIterate(t *testing.T) {
	v := Of(1, 2, 3, 4)
	i := 0
	for it := v.Iter(); it.HasNext(); {
		n := it.Next()
		expected, _ := v.Get(i)
		assert.Equal(t, expected, n)
		i++
	}
}

func Test_Vector_Sort(t *testing.T) {
	var (
		elements = Of(3, 1, 5, 6, 8, 0, -2)
		sorted   = elements.Sort(func(e1, e2 int) bool { return e1 < e2 })
	)
	assert.Equal(t, Of(-2, 0, 1, 3, 5, 6, 8), sorted)
}

func Test_Vector_SortStructByField(t *testing.T) {
	var (
		anonymous = &user{"Anonymous", 0}
		cherlie   = &user{"Cherlie", 25}
		alise     = &user{"Alise", 20}
		bob       = &user{"Bob", 19}

		elements     = Of(anonymous, cherlie, alise, bob)
		sortedByName = Sort(elements, (*user).Name)
		sortedByAge  = Sort(elements, (*user).Age)
	)
	assert.Equal(t, Of(alise, anonymous, bob, cherlie), sortedByName)
	assert.Equal(t, Of(anonymous, bob, alise, cherlie), sortedByAge)
}

type user struct {
	name string
	age  int
}

func (u *user) Name() string { return u.name }
func (u *user) Age() int     { return u.age }
