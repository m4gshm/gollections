package test

import (
	"testing"

	oMapIter "github.com/m4gshm/gollections/immutable/ordered/map_/iter"
	"github.com/m4gshm/gollections/map_/iter"
)

func Test_Key_Zero_Safety(t *testing.T) {
	var it iter.KeyIter[int, string]

	it.Next()
	it.Cap()

	//OrderedEmbedMapKVIter
}

func Test_OrderedEmbedMapKVIter_Safety(t *testing.T) {
	var it oMapIter.OrderedMapIter[int, string]

	it.Next()
	it.Cap()
}
