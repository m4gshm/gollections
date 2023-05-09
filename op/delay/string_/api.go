package string_

import (
	"strings"
)

// Of returns spring builder
func Of(parts ...string) func() string {
	return func() string { return strings.Join(parts, "") }
}
