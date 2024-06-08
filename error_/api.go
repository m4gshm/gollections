// Package error_ provides generic functions for error handling.
package error_

import (
	"errors"
)

// As finds the first error in err's tree that matches E type and retunrs the one and ok==true.
// If nothing to found it will return ok==false.
func As[E error](err error) (out E, ok bool) {
	return out, errors.As(err, &out)
}
