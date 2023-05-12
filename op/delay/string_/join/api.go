package join

import "github.com/m4gshm/gollections/op/delay/string_"

func NonEmpty(joiner string) func(first, second string) string {
	return string_.JoinNonEmpty(joiner)
}
