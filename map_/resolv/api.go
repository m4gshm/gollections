package resolv

import (
	"github.com/m4gshm/gollections/op"
)

// FirstVal - ToMap value resolver
func FirstVal[K, V any](exists bool, key K, old, new V) V { return op.IfElse(exists, old, new) }

// LastVal - ToMap value resolver
func LastVal[K, V any](exists bool, key K, old, new V) V { return new }
