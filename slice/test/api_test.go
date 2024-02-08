package test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	_less "github.com/m4gshm/gollections/break/predicate/less"
	_more "github.com/m4gshm/gollections/break/predicate/more"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/delay/chain"
	"github.com/m4gshm/gollections/op/delay/string_/wrap"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/clone/reverse"
	csort "github.com/m4gshm/gollections/slice/clone/sort"
	cstablesort "github.com/m4gshm/gollections/slice/clone/stablesort"
	"github.com/m4gshm/gollections/slice/conv"
	"github.com/m4gshm/gollections/slice/convert"
	"github.com/m4gshm/gollections/slice/filter"
	"github.com/m4gshm/gollections/slice/first"
	"github.com/m4gshm/gollections/slice/last"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/m4gshm/gollections/slice/split"
	"github.com/m4gshm/gollections/slice/stablesort"
)

func Test_Range(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2), range_.Of(-1, 3))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), range_.Of(3, -2))
	assert.Equal(t, slice.Of(1), range_.Of(1, 2))
	assert.Nil(t, range_.Of(1, 1))
}

func Test_RangeClose(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), range_.Closed(-1, 3))
	assert.Equal(t, slice.Of(3, 2, 1, 0, -1), range_.Closed(3, -1))
	assert.Equal(t, slice.Of(1), range_.Closed(1, 1))
}

func Test_Reverse(t *testing.T) {
	src := range_.Closed(3, -1)
	reversed := slice.Reverse(src)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), reversed)
	assert.Equal(t, (*reflect.SliceHeader)(unsafe.Pointer(&src)).Data, (*reflect.SliceHeader)(unsafe.Pointer(&reversed)).Data)
}

func Test_ReverseCloned(t *testing.T) {
	src := range_.Closed(3, -1)
	reversed := reverse.Of(src)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), reversed)
	assert.NotEqual(t, (*reflect.SliceHeader)(unsafe.Pointer(&src)).Data, (*reflect.SliceHeader)(unsafe.Pointer(&reversed)).Data)
}

func Test_Clone(t *testing.T) {
	type entity struct{ val string }
	var (
		first  = entity{"first"}
		second = entity{"second"}
		third  = entity{"third"}

		entities = []*entity{&first, &second, &third}
		c        = clone.Of(entities)
	)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Same(t, entities[i], c[i])
	}
}

func Test_DeepClone(t *testing.T) {
	type entity struct{ val string }
	var (
		first  = entity{"first"}
		second = entity{"second"}
		third  = entity{"third"}

		entities = []*entity{&first, &second, &third}
		c        = clone.Deep(entities, clone.Ptr[entity])
	)

	assert.Equal(t, entities, c)
	assert.NotSame(t, entities, c)

	for i := range entities {
		assert.Equal(t, entities[i], c[i])
		assert.NotSame(t, entities[i], c[i])
	}
}

func Test_PointerClone(t *testing.T) {

	s1 := "test"
	p1 := &s1

	p2 := clone.Ptr(p1)

	assert.Equal(t, s1, *p2)
	assert.NotSame(t, p1, p2)
}

func Test_ReduceSum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Reduce(s, op.Sum[int])
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_ConvertAndReduce(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := convert.AndReduce(s, func(i int) int { return i * i }, op.Sum[int])
	assert.Equal(t, 1+3*3+5*5+7*7+9*9+11*11, r)
}

func Test_ConvAndReduce(t *testing.T) {
	s := slice.Of("1", "3", "5", "7", "9", "11")
	r, err := conv.AndReduce(s, strconv.Atoi, op.Sum[int])
	assert.NoError(t, err)
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_Sum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Sum(s)
	assert.Equal(t, 1+3+5+7+9+11, r)
}

func Test_First(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, ok := first.Of(s, func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := slice.First(s, func(i int) bool { return i > 12 })
	assert.False(t, nook)
}

func Test_Firstt(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, ok, _ := slice.Firstt(s, _more.Than(5))
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook, _ := slice.Firstt(s, _more.Than(12))
	assert.False(t, nook)

	_, _, err := slice.Firstt(s, func(i int) (bool, error) { return true, errors.New("abort") })
	assert.Error(t, err)
}

func Test_Last(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, ok := last.Of(s, func(i int) bool { return i < 9 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook := slice.Last(s, func(i int) bool { return i < 1 })
	assert.False(t, nook)
}

func Test_Lastt(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, ok, _ := slice.Lastt(s, _less.Than(9))
	assert.True(t, ok)
	assert.Equal(t, 7, r)

	_, nook, _ := slice.Lastt(s, _less.Than(1))
	assert.False(t, nook)

	_, _, err := slice.Lastt(s, func(i int) (bool, error) { return true, errors.New("abort") })
	assert.Error(t, err)
}

var absPath = op.IfElse(runtime.GOOS == "windows", "c:\\home\\user", "/home/user")
var absPath2 = op.IfElse(runtime.GOOS == "windows", "c:\\usr\\bin", "/usr/bin")

func TestConv(t *testing.T) {
	if homeDir, err := os.UserHomeDir(); err != nil {
		t.Error(err)
	} else if err := os.Chdir(homeDir); err != nil {
		t.Error(err)
	} else if abs, err := slice.Conv(slice.Of(absPath, "././inTemp"), filepath.Abs); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, slice.Of(absPath, filepath.Join(homeDir, "inTemp")), abs)
	}
}

func Test_Convert(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Convert(s, strconv.Itoa)
	assert.Equal(t, []string{"1", "3", "5", "7", "9", "11"}, r)
}

func Test_Conv(t *testing.T) {
	s := slice.Of("1", "3", "5", "7", "_9", "11")
	r, err := slice.Conv(s, strconv.Atoi)
	var expected *strconv.NumError
	assert.ErrorAs(t, err, &expected)
	assert.Equal(t, []int{1, 3, 5, 7}, r)
}

func Test_ConvertWithIndex(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.ConvertIndexed(s, func(index int, elem int) int { return index + elem })
	assert.Equal(t, slice.Of(1, 1+3, 2+5, 3+7, 4+9, 5+11), r)
}

func Test_ConvertNotNil(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = []*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}
		result   = convert.NotNil(source, func(e *entity) string { return e.val })
		expected = []string{"first", "third", "fifth"}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertValues(t *testing.T) {
	type entity struct{ val string }
	var (
		source   = []*entity{{"first"}, nil, {"third"}, nil, {"fifth"}}
		result   = convert.NotNil(source, as.Val[entity])
		expected = []entity{{"first"}, {"third"}, {"fifth"}}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertToNotNil(t *testing.T) {
	type entity struct{ val *string }
	var (
		first    = "first"
		third    = "third"
		fifth    = "fifth"
		source   = []entity{{&first}, {}, {&third}, {}, {&fifth}}
		result   = convert.ToNotNil(source, func(e entity) *string { return e.val })
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, result)
}

func Test_ConvertNilSafe(t *testing.T) {
	type entity struct{ val *string }
	var (
		first    = "first"
		third    = "third"
		fifth    = "fifth"
		source   = []*entity{{&first}, {}, {&third}, nil, {&fifth}}
		result   = convert.NilSafe(source, func(e *entity) *string { return e.val })
		expected = []*string{&first, &third, &fifth}
	)
	assert.Equal(t, expected, result)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.FilterAndConvert(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndex(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.FilterAndConvertIndexed(s, func(_ int, elem int) bool { return even(elem) }, func(index int, elem int) string { return strconv.Itoa(index + elem) })
	assert.Equal(t, []string{"6", "13"}, r)
}

func Test_ConvertFilteredInplace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.ConvertCheck(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.ConvertCheckIndexed(s, func(index int, elem int) (string, bool) { return strconv.Itoa(index + elem), even(elem) })
	assert.Equal(t, []string{"6", "13"}, r)
}

func Test_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flat(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Flat(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f, err := slice.Flatt(md, func(i []int) ([]int, error) { return i, op.IfElse(len(i) == 2, errors.New("abort"), nil) })
	assert.Error(t, err)
	e := []int{1, 2, 3, 4}
	assert.Equal(t, e, f)
}

func Benchmark_Flatt(b *testing.B) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}

	for i := 0; i < b.N; i++ {
		_ = slice.Flat(md, func(i []int) []int { return i })
	}
}

func Benchmark_Flatt_Convert_AsIs(b *testing.B) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}

	for i := 0; i < b.N; i++ {
		_ = slice.FlatAndConvert(md, func(i []int) []int { return i }, as.Is[int])
	}
}

func Test_FlattFilter(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FilterAndFlat(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i })
	e := []int{1, 2, 3, 5, 6}
	assert.Equal(t, e, f)
}

func Test_FlattElemFilter(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FlatAndFiler(md, func(i []int) []int { return i }, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, f)
}

func Test_FilterAndFlattFilt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FilterFlatFilter(md, func(from []int) bool { return len(from) > 1 }, func(i []int) []int { return i }, even)
	e := []int{2, 6}
	assert.Equal(t, e, f)
}

func Test_Filter(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.Filter(s, even)
	assert.Equal(t, slice.Of(4, 8), r)
}

func Test_FilterConvertFilter(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := filter.ConvertFilter(s, even, func(i int) int { return i * 2 }, even)
	assert.Equal(t, slice.Of(8, 16), r)
}

func Test_Filt(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r, err := slice.Filt(s, func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) })
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4, 8), r)
}

func Test_Filt2(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r, err := slice.Filt(s, func(i int) (bool, error) {
		ok := i <= 7
		return ok && even(i), op.IfElse(ok, nil, errors.New("abort"))

	})
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)
}

func Test_FiltAndConv(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r, _ := slice.FiltAndConv(s, func(v int) (bool, error) { return v%2 == 0, nil }, func(i int) (int, error) { return i * 2, nil })
	assert.Equal(t, slice.Of(8, 16), r)
}

func Test_StringRepresentation(t *testing.T) {
	var (
		expected = fmt.Sprint(slice.Of(1, 2, 3, 4))
		actual   = slice.ToString(slice.Of(1, 2, 3, 4))
	)
	assert.Equal(t, expected, actual)
}

func Test_StringReferencesRepresentation(t *testing.T) {
	var (
		expected       = fmt.Sprint(slice.Of(1, 2, 3, 4))
		i1, i2, i3, i4 = 1, 2, 3, 4
		actual         = slice.ToStringRefs(slice.Of(&i1, &i2, &i3, &i4))
	)
	assert.Equal(t, expected, actual)
}

func Test_Filtering(t *testing.T) {
	r := slice.Filter([]int{1, 2, 3, 4, 5, 6}, func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, r)
}

func Test_BehaveAsStrings(t *testing.T) {
	type TypeBasedOnString string
	type ArrayTypeBasedOnString []TypeBasedOnString

	vals := ArrayTypeBasedOnString{"1", "2", "3"}
	strs := slice.BehaveAsStrings(vals)

	assert.Equal(t, []string{"1", "2", "3"}, strs)
	pvals := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data)
	pstrs := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&strs)).Data)
	assert.Equal(t, pvals, pstrs)
}

func Test_StringsBehaveAs(t *testing.T) {
	type TypeBasedOnString string
	type ArrayTypeBasedOnString []TypeBasedOnString

	vals := []string{"1", "2", "3"}
	strs := slice.StringsBehaveAs[ArrayTypeBasedOnString](vals)

	assert.Equal(t, ArrayTypeBasedOnString{"1", "2", "3"}, strs)
	pvals := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data)
	pstrs := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&strs)).Data)
	assert.Equal(t, pvals, pstrs)
}

func Test_StringsBehaveAs2(t *testing.T) {
	type TypeBasedOnString string

	vals := []string{"1", "2", "3"}
	strs := slice.StringsBehaveAs[[]TypeBasedOnString](vals)

	assert.Equal(t, []TypeBasedOnString{"1", "2", "3"}, strs)
	pvals := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data)
	pstrs := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&strs)).Data)
	assert.Equal(t, pvals, pstrs)
}

type rows[T any] struct {
	row    []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.row) }
func (r *rows[T]) next() (T, error) { e := r.row[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := slice.OfLoop(stream, (*rows[int]).hasNext, (*rows[int]).next)

	assert.Equal(t, slice.Of(1, 2, 3), result)
}

func Test_Sort(t *testing.T) {
	src := range_.Closed(3, -1)
	sorted := sort.Asc(src)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), sorted)
	assert.Equal(t, (*reflect.SliceHeader)(unsafe.Pointer(&src)).Data, (*reflect.SliceHeader)(unsafe.Pointer(&sorted)).Data)
}

func Test_SortCloned(t *testing.T) {
	src := range_.Closed(3, -1)
	sorted := csort.Asc(src)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), sorted)
	assert.NotEqual(t, (*reflect.SliceHeader)(unsafe.Pointer(&src)).Data, (*reflect.SliceHeader)(unsafe.Pointer(&sorted)).Data)
}

func Test_StableSort(t *testing.T) {
	src := range_.Closed(3, -1)
	sorted := stablesort.Asc(src)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), sorted)
	assert.Equal(t, (*reflect.SliceHeader)(unsafe.Pointer(&src)).Data, (*reflect.SliceHeader)(unsafe.Pointer(&sorted)).Data)
}

func Test_StableSortCloned(t *testing.T) {
	src := range_.Closed(3, -1)
	sorted := cstablesort.Asc(src)
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), sorted)
	assert.NotEqual(t, (*reflect.SliceHeader)(unsafe.Pointer(&src)).Data, (*reflect.SliceHeader)(unsafe.Pointer(&sorted)).Data)
}

func Test_MatchAny(t *testing.T) {
	elements := slice.Of(1, 2, 3, 4)

	ok := slice.HasAny(elements, eq.To(4))
	assert.True(t, ok)

	noOk := slice.HasAny(elements, more.Than(5))
	assert.False(t, noOk)
}

func Test_Empty(t *testing.T) {
	assert.False(t, slice.Empty(slice.Of(1)))
	assert.True(t, slice.Empty(slice.Of[int]()))
	assert.True(t, slice.Empty[[]int](nil))
}

func Test_SplitTwo(t *testing.T) {
	first, second := slice.SplitTwo(slice.Of("1a", "2b", "3c"), func(s string) (string, string) { return string(s[0]), string(s[1]) })

	assert.Equal(t, slice.Of("1", "2", "3"), first)
	assert.Equal(t, slice.Of("a", "b", "c"), second)
}

func Test_SplitTwo2(t *testing.T) {
	byIndex := func(i int) func(string) string { return func(s string) string { return string(s[i]) } }

	first, second := split.Of(slice.Of("1a", "2b", "3c"), byIndex(0), byIndex(1))

	assert.Equal(t, slice.Of("1", "2", "3"), first)
	assert.Equal(t, slice.Of("a", "b", "c"), second)
}

func Test_SplitAndReduce(t *testing.T) {
	byIndex := func(i int) func(string) string { return func(s string) string { return string(s[i]) } }

	first, second := split.AndReduce(slice.Of("1a", "2b", "3c"), byIndex(0), chain.Of(byIndex(1), wrap.By("{", "}")), op.Sum[string], op.Sum[string])

	assert.Equal(t, "123", first)
	assert.Equal(t, "{a}{b}{c}", second)
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := slice.OfIndexed(len(indexed), func(i int) string { return indexed[i] })
	assert.Equal(t, indexed, result)
}

func Test_PeekWhile(t *testing.T) {
	expected := slice.Of(1, 3, 5, 7, 9, 11)

	s := []int{}

	slice.PeekWhile(expected, func(e int) bool {
		s = append(s, e)
		return true
	})

	assert.Equal(t, expected, s)
}
