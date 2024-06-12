//go:build goexperiment.rangefunc

package test

import (
	"strconv"
	"testing"

	"github.com/m4gshm/gollections/collection/immutable/ordered/map_"
	"github.com/m4gshm/gollections/k"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/stretchr/testify/assert"
)

var testMap = map_.Of(k.V(1, "10"), k.V(2, "20"), k.V(3, "30"), k.V(5, "50"), k.V(7, "70"), k.V(8, "80"), k.V(9, "90"), k.V(11, "110"))

func Test_AllFiltered(t *testing.T) {
	s := []string{}

	for _, v := range seq2.Filter(testMap.All, func(k int, v string) bool { return k%2 == 0 }) {
		s = append(s, v)
	}

	assert.Equal(t, slice.Of("20", "80"), sort.Asc(s))
}

func Test_AllConverted(t *testing.T) {
	i := []int{}

	for _, e := range seq2.Convert(testMap.All, func(k int, v string) int { c, _ := strconv.Atoi(v); return c }) {
		i = append(i, e)
	}

	assert.Equal(t, slice.Of(10, 20, 30, 50, 70, 80, 90, 110), i)
}
