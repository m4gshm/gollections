package breakableloop

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/break/loop"
	"github.com/stretchr/testify/assert"
)

type (
	next[T any]      func() (element T, ok bool, err error)
	kvNext[K, V any] func() (key K, value V, ok bool, err error)
)

func Test_Slice_Vs_Loop(t *testing.T) {

	iter := loop.Conv(loop.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
	result, err := loop.Slice(iter.Next)

	assert.Equal(t, []int{1, 2, 3}, result)
	assert.ErrorContains(t, err, "invalid syntax")
}

// func Test_Loop_Over_Sql_Rows(t *testing.T) {
// 	type User struct {
// 		name string
// 		age  int
// 	}
// 	var rows *sql.Rows

// 	users := loop.New(rows, (*sql.Rows).Next, func(row *sql.Rows) (*User, error) {
// 		user := User{}
// 		return &user, rows.Scan(&user.name, &user.age)
// 	}).Filter(func(u *User) bool { return u.age > 18 })

// }
