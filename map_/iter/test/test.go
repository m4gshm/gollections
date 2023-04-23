package test

import (
	"testing"

	"github.com/m4gshm/gollections/map_/iter"
)

func Test_Key_Zero_Safety(t *testing.T) {
	var it iter.Key[int, string]

	it.Next()
	it.Cap()

	//OrderedEmbedMapKVIter
}

func Test_OrderedEmbedMapKVIter_Safety(t *testing.T) {
	var it iter.OrderedEmbedMapKVIter[int, string]

	it.Next()
	it.Cap()
}
