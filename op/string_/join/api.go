package join

import "github.com/m4gshm/gollections/op/string_"

func NoEmpty(first, joiner, second string) string {
	return string_.JoinNoEmpty(first, joiner, second)
}