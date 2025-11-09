package seqexamples

import (
	"testing"

	"github.com/m4gshm/gollections/convert/as"
	sql "github.com/m4gshm/gollections/internal/examples/seqexamples/sql"
	"github.com/m4gshm/gollections/seqe"
	"github.com/stretchr/testify/assert"
)

func Test_rows_OfNext(t *testing.T) {

	var rows sql.Rows = selectUsers()

	getUser := func(u *User) error { return rows.Scan(&u.name, &u.age) }
	rowSeq := seqe.OfNext(rows.Next, getUser)
	usersByAge, err := seqe.Group(rowSeq, User.Age, as.Is)

	assert.Len(t, usersByAge, 1)
	assert.Equal(t, "Alice", usersByAge[21][0].Name())
	assert.NoError(t, err)
}

func Test_rows_plain_old(t *testing.T) {

	var rows sql.Rows = selectUsers()

	var usersByAge = map[int][]User{}
	var err error
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.name, &u.age); err != nil {
			break
		}
		usersByAge[u.age] = append(usersByAge[u.age], u)
	}

	assert.Len(t, usersByAge, 1)
	assert.Equal(t, "Alice", usersByAge[21][0].Name())
	assert.NoError(t, err)
}

func Test_rows_OfNextGet(t *testing.T) {

	var rows sql.Rows = selectUsers()

	refUsersByAge, err := seqe.Group(seqe.OfNextGet(rows.Next, func() (*User, error) {
		var u User
		return &u, rows.Scan(&u.name, &u.age)
	}), (*User).Age, as.Is)

	assert.Len(t, refUsersByAge, 1)
	assert.Equal(t, "Alice", refUsersByAge[21][0].Name())
	assert.NoError(t, err)
}

func selectUsers() sql.Rows {
	return sql.Rows{}
}
