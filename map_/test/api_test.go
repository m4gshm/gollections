package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/clone"
	"github.com/m4gshm/gollections/map_/group"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone/sort"
)

type entity struct{ val string }

var (
	first  = entity{"1_first"}
	second = entity{"2_second"}
	third  = entity{"3_third"}

	entities = map[int]*entity{1: &first, 2: &second, 3: &third}
)

func Test_Clone(t *testing.T) {
	c := clone.Of(entities)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for k := range entities {
		assert.Same(t, entities[k], c[k])
	}
}

func Test_DeepClone(t *testing.T) {
	c := clone.Deep(entities, func(e *entity) *entity { return ptr.Of(*e) })

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Equal(t, entities[i], c[i])
		assert.NotSame(t, entities[i], c[i])
	}
}

func Test_Keys(t *testing.T) {
	keys := map_.Keys(entities)
	assert.Equal(t, slice.Of(1, 2, 3), sort.Of(keys))
}

func Test_Values(t *testing.T) {
	values := map_.Values(entities)
	assert.Equal(t, slice.Of(&first, &second, &third), sort.By(values, func(e *entity) string { return e.val }))
}

func Test_ConvertValues(t *testing.T) {
	var strValues map[int]string = map_.ConvertValues(entities, func(e *entity) string { return e.val })

	assert.Equal(t, "1_first", strValues[1])
	assert.Equal(t, "2_second", strValues[2])
	assert.Equal(t, "3_third", strValues[3])
}

func Test_ValuesConverted(t *testing.T) {
	var values []string = map_.ValuesConverted(entities, func(e *entity) string { return e.val })
	assert.Equal(t, slice.Of("1_first", "2_second", "3_third"), sort.Of(values))
}

type rows[T any] struct {
	in     []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.in) }
func (r *rows[T]) next() (T, error) { e := r.in[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := map_.OfLoop(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	})

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_OfLoopResolv(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3, 4), 0}
	result, _ := map_.OfLoopResolv(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	}, resolv.Last[bool, int])

	assert.Equal(t, 4, result[true])
	assert.Equal(t, 3, result[false])
}

func Test_Generate(t *testing.T) {
	counter := 0
	result, _ := map_.Generate(func() (bool, int, bool, error) { counter++; return counter%2 == 0, counter, counter < 4, nil })

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_GenerateResolv(t *testing.T) {
	counter := 0
	result, _ := map_.GenerateResolv(func() (bool, int, bool, error) {
		counter++
		return counter%2 == 0, counter, counter < 5, nil
	}, func(exists bool, k bool, old, new int) int {
		return op.IfElse(exists, op.IfElse(k, new, old), new)
	})

	assert.Equal(t, 4, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_GroupOfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := group.OfLoop(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	})

	assert.Equal(t, slice.Of(2), result[true])
	assert.Equal(t, slice.Of(1, 3), result[false])
}

func Test_StringRepresentation(t *testing.T) {
	order := slice.Of(4, 3, 2, 1)
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}
	actual := map_.ToStringOrdered(order, elements)

	expected := "[4:4 3:3 2:2 1:1]"
	assert.Equal(t, expected, actual)
}

func Test_Reduce(t *testing.T) {
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}
	k, _ := map_.Reduce(elements, func(k, k2 int, v, v2 string) (int, string) {
		return k + k2, ""
	})

	assert.Equal(t, 1+2+3+4, k)
}

func Test_MatchAny(t *testing.T) {
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}
	ok := map_.HasAny(elements, func(k int, v string) bool {
		return k == 2 || v == "4"
	})

	assert.True(t, ok)

	noOk := map_.HasAny(elements, func(k int, v string) bool {
		return k > 5
	})

	assert.False(t, noOk)
}

func Test_ToSlice(t *testing.T) {
	result := map_.ToSlice(entities, func(key int, val *entity) string { return strconv.Itoa(key) + ":" + val.val })
	assert.Equal(t, slice.Of("1:1_first", "2:2_second", "3:3_third"), sort.Of(result))
}

func Test_ToSliceErrorable(t *testing.T) {
	result, _ := map_.ToSlicee(entities, func(key int, val *entity) (int, error) {
		v, err := strconv.Atoi(string(val.val[0]))
		return v + key, err
	})
	assert.Equal(t, slice.Of(2, 4, 6), sort.Of(result))
}

func Test_Filter(t *testing.T) {
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}

	result := map_.Filter(elements, func(key int, val string) bool { return key <= 2 || val == "4" })
	check := map_.KeyChecker(result)
	assert.Equal(t, 3, len(result))
	assert.True(t, check(1))
	assert.True(t, check(2))
	assert.False(t, check(3))
	assert.True(t, check(4))
}

func Test_FilterKeys(t *testing.T) {
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}

	result := map_.FilterKeys(elements, func(key int) bool { return key <= 2 })
	check := map_.KeyChecker(result)
	assert.Equal(t, 2, len(result))
	assert.True(t, check(1))
	assert.True(t, check(2))
	assert.False(t, check(3))
	assert.False(t, check(4))
}

func Test_FilterValues(t *testing.T) {
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}

	result := map_.FilterValues(elements, func(val string) bool { return val <= "2" })
	check := map_.KeyChecker(result)
	assert.Equal(t, 2, len(result))
	assert.True(t, check(1))
	assert.True(t, check(2))
	assert.False(t, check(3))
	assert.False(t, check(4))
}

func Test_Getter(t *testing.T) {
	elements := map[int]string{4: "4", 2: "2", 1: "1", 3: "3"}

	getter := map_.Getter(elements)
	assert.Equal(t, "1", getter(1))
	assert.Equal(t, "2", getter(2))
	assert.Equal(t, "3", getter(3))
	assert.Equal(t, "4", getter(4))

	nilGetter := map_.Getter[map[int]string](nil)

	assert.Equal(t, "", nilGetter(0))

	getterOk := map_.GetterOk(elements)

	_, ok := getterOk(1)
	assert.True(t, ok)
	_, ok = getterOk(0)
	assert.False(t, ok)
}
