package join

import "github.com/m4gshm/gollections/op/delay/string_"

func JoinNoEmpty(joiner string) func(first, second string) string {
	return string_.JoinNoEmpty(joiner)
}
