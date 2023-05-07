package string_

import (
	"strings"
)

func Of(parts ...string) func() string {
	return func() string { return strings.Join(parts, "") }
	// return sum.Of(parts...)
	// return func() string { return curry.Second(strings.Join, "")(parts) }
}
