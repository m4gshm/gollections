package test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	_less "github.com/m4gshm/gollections/break/predicate/less"
	_more "github.com/m4gshm/gollections/break/predicate/more"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/delay/chain"
	"github.com/m4gshm/gollections/op/delay/string_/wrap"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/more"
	"github.com/m4gshm/gollections/seq"
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
	assert.NotSame(t, &entities, &c)

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
	assert.NotSame(t, &entities, &c)

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

func Test_ReduceeSum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, _ := slice.Reducee(s, func(i1, i2 int) (int, error) { return i1 + i2, nil })
	assert.Equal(t, 1+3+5+7+9+11, r)

	r2, err := slice.Reducee(s, func(i1, i2 int) (int, error) { return i1 + i2, op.IfElse(i2 == 7, errors.New("abort"), nil) })
	assert.Error(t, err)
	assert.Equal(t, 1+3+5+7, r2)
}

func Test_AccumSum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Accum(100, s, op.Sum[int])
	assert.Equal(t, 100+1+3+5+7+9+11, r)
}

func Test_AccummSum(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r, _ := slice.Accumm(100, s, func(i1, i2 int) (int, error) { return i1 + i2, nil })
	assert.Equal(t, 100+1+3+5+7+9+11, r)
	r2, err := slice.Accumm(100, s, func(i1, i2 int) (int, error) { return i1 + i2, op.IfElse(i2 == 7, errors.New("abort"), nil) })
	assert.Error(t, err)
	assert.Equal(t, 100+1+3+5+7, r2)
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

	s = slice.Of("1", "3", "5", "_7", "9", "11")

	r, err = conv.AndReduce(s, strconv.Atoi, op.Sum[int])
	assert.ErrorContains(t, err, "parsing \"_7\": invalid syntax")
	assert.Equal(t, 1+3+5, r)

	s = slice.Of("_1")

	r, err = conv.AndReduce(s, strconv.Atoi, op.Sum[int])
	assert.ErrorContains(t, err, "parsing \"_1\": invalid syntax")
	assert.Equal(t, 0, r)
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

	_, _, err := slice.Firstt(s, func(_ int) (bool, error) { return true, errors.New("abort") })
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

	_, _, err := slice.Lastt(s, func(_ int) (bool, error) { return true, errors.New("abort") })
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

func Test_ConvOK(t *testing.T) {
	s := slice.Of("1", "3", "4", "5", "7", "8", "_9", "11", "12")
	r, err := slice.ConvOK(s, func(v string) (int, bool, error) {
		i, err := strconv.Atoi(v)
		if err != nil {
			return i, false, err
		}
		return i, even(i), nil

	})
	var expected *strconv.NumError
	assert.ErrorAs(t, err, &expected)
	assert.Equal(t, []int{4, 8}, r)

	s = slice.Of("1", "3", "4", "5", "7", "8", "9", "11", "12")
	r, err = slice.ConvOK(s, func(v string) (int, bool, error) {
		i, err := strconv.Atoi(v)
		return i, true, err

	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3, 4, 5, 7, 8, 9, 11, 12}, r)

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
		source   = []*entity{{"first"}, nil, {"third"}, {}, {"fifth"}}
		result   = convert.NotNil(source, as.Val[entity])
		expected = []entity{{"first"}, {"third"}, {}, {"fifth"}}
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

func Test_ConvertOK(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.ConvertOK(s, func(i int) (string, bool) { return strconv.Itoa(i), even(i) })
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.ConvertCheckIndexed(s, func(index int, elem int) (string, bool) { return strconv.Itoa(index + elem), even(elem) })
	assert.Equal(t, []string{"6", "13"}, r)
}

func Test_Flat(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flat(md, as.Is)
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_FlatSeq(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.FlatSeq(md, slices.Values)
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f, err := slice.Flatt(md, func(i []int) ([]int, error) { return i, op.IfElse(len(i) == 2, errors.New("abort"), nil) })
	assert.Error(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, f)

	f, err = slice.Flatt(md, func(i []int) ([]int, error) { return i, nil })
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, f)
}

func Test_FlattSeq(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	transform := func(i int) (int, error) {
		return i, op.IfElse(i == 5, errors.New("abort"), nil)
	}
	f, err := slice.FlattSeq(md, func(i []int) seq.SeqE[int] { return seq.Conv(seq.Of(i...), transform) })
	assert.Error(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, f)

	f, err = slice.FlattSeq(md, func(i []int) seq.SeqE[int] {
		return seq.SeqE[int](seq.ToSeq2(seq.Of(i...), func(i int) (int, error) { return i, nil }))
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, f)

}

func Test_FlatAndConvert(t *testing.T) {
	md := slice.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...)
	f := slice.FlatAndConvert(md, func(i []int) []int { return i }, strconv.Itoa)
	e := []string{"1", "2", "3", "4", "5", "6"}
	assert.Equal(t, e, f)
}

func Benchmark_Flat(b *testing.B) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}

	for i := 0; i < b.N; i++ {
		_ = slice.Flat(md, as.Is)
	}
}

func Benchmark_Flat_Convert_AsIs(b *testing.B) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}

	for i := 0; i < b.N; i++ {
		_ = slice.FlatAndConvert(md, as.Is, as.Is)
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
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 12)
	r, err := slice.Filt(s, func(i int) (bool, error) { return even(i), op.IfElse(i > 7, errors.New("abort"), nil) })
	assert.Error(t, err)
	assert.Equal(t, slice.Of(4), r)

	r, err = slice.Filt(s, func(i int) (bool, error) { return even(i), nil })
	assert.NoError(t, err)
	assert.Equal(t, slice.Of(4, 8, 12), r)
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

func Test_AppendFiltAndConv(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r, err := slice.AppendFiltAndConv(s, []int{}, func(v int) (bool, error) { return v%2 == 0, nil }, func(i int) (int, error) { return i * 2, nil })
	assert.Equal(t, slice.Of(8, 16), r)
	assert.NoError(t, err)

	r, err = slice.AppendFiltAndConv(s, []int{}, func(v int) (bool, error) { return v%2 == 0, op.IfElse(v == 9, errors.New("abort"), nil) }, func(i int) (int, error) { return i * 2, nil })

	assert.Equal(t, slice.Of(8, 16), r)
	assert.ErrorContains(t, err, "abort")

	r, err = slice.AppendFiltAndConv(s, []int{}, func(v int) (bool, error) { return v%2 == 0, nil }, func(i int) (int, error) { return i * 2, op.IfElse(i == 8, errors.New("abort"), nil) })

	assert.Equal(t, slice.Of(8), r)
	assert.ErrorContains(t, err, "abort")
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

type Rows[T any] struct {
	row    []T
	cursor int
}

func (r *Rows[T]) Reset()             { r.cursor = 0 }
func (r *Rows[T]) Next() bool         { return r.cursor < len(r.row) }
func (r *Rows[T]) Scan(dest *T) error { *dest = r.row[r.cursor]; r.cursor++; return nil }

func Test_OffNextPush(t *testing.T) {
	var (
		rows        = &Rows[int]{slice.Of(1, 2, 3), 0}
		result, err = slice.OfNext(rows.Next, rows.Scan)
		expected    = slice.Of(1, 2, 3)
	)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)

	rows.Reset()

	result, err = slice.OfSourceNext(rows, (*Rows[int]).Next, func(r *Rows[int], our *int) error {
		return r.Scan(our)
	})

	assert.Equal(t, expected, result)
	assert.NoError(t, err)
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
	assert.False(t, slice.IsEmpty(slice.Of(1)))
	assert.True(t, slice.IsEmpty(slice.Of[int]()))
	assert.True(t, slice.IsEmpty[[]int](nil))
}

func Test_SplitTwo(t *testing.T) {
	first, second := slice.SplitTwo(slice.Of("1a", "2b", "3c"), func(s string) (string, string) { return string(s[0]), string(s[1]) })

	assert.Equal(t, slice.Of("1", "2", "3"), first)
	assert.Equal(t, slice.Of("a", "b", "c"), second)
}

func Test_SplitThree(t *testing.T) {
	first, second, third := slice.SplitThree(slice.Of("1a#", "2b$", "3c%"), func(s string) (string, string, string) { return string(s[0]), string(s[1]), string(s[2]) })

	assert.Equal(t, slice.Of("1", "2", "3"), first)
	assert.Equal(t, slice.Of("a", "b", "c"), second)
	assert.Equal(t, slice.Of("#", "$", "%"), third)
}

func Test_SplitTwo2(t *testing.T) {
	byIndex := func(i int) func(string) string { return func(s string) string { return string(s[i]) } }

	first, second := split.Of(slice.Of("1a", "2b", "3c"), byIndex(0), byIndex(1))

	assert.Equal(t, slice.Of("1", "2", "3"), first)
	assert.Equal(t, slice.Of("a", "b", "c"), second)
}

func Test_SplitAndReduce(t *testing.T) {
	byIndex := func(i int) func(string) string { return func(s string) string { return string(s[i]) } }

	first, second := split.AndReduce(slice.Of("1a", "2b", "3c"), byIndex(0), chain.Of(byIndex(1), wrap.By("{", "}")), op.Sum, op.Sum)

	assert.Equal(t, "123", first)
	assert.Equal(t, "{a}{b}{c}", second)
}

func Test_OfIndexed(t *testing.T) {
	indexed := slice.Of("0", "1", "2", "3", "4")
	result := slice.OfIndexed(len(indexed), func(i int) string { return indexed[i] })
	assert.Equal(t, indexed, result)
}

func Test_Series(t *testing.T) {
	assert.Equal(t, slice.Of(-1, 0, 1, 2, 3), slice.Series(-1, func(prev int) (int, bool) { return prev + 1, prev < 3 }))
}

func Test_PeekWhile(t *testing.T) {
	expected := slice.Of(1, 3, 5, 7, 9, 11)

	s := []int{}

	slice.WalkWhile(expected, func(e int) bool {
		s = append(s, e)
		return true
	})

	assert.Equal(t, expected, s)
}

func Test_Slice_ToMapResolvOrder(t *testing.T) {
	var (
		even          = func(v int) bool { return v%2 == 0 }
		order, groups = slice.MapResolvOrder(slice.Of(2, 1, 1, 2, 4, 3, 1), even, as.Is[int], resolv.Slice)
	)
	assert.Equal(t, []int{1, 1, 3, 1}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.Equal(t, []bool{true, false}, order)
}

func Test_Slice_AppendMapResolv(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = slice.AppendMapResolv(slice.Of(2, 1, 1, 2, 4, 3, 1), even, as.Is[int], resolv.Slice, nil)
	)
	assert.Equal(t, []int{1, 1, 3, 1}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
}

func Test_Slice_AppendMapResolvv(t *testing.T) {
	even := func(v int) bool { return v%2 == 0 }
	groups, err := slice.AppendMapResolvv(slice.Of(2, 1, 1, 2, 4, 3, 1), func(i int) (bool, int, error) {
		return even(i), i, nil
	}, func(_ bool, _ bool, e []int, v int) ([]int, error) {
		return append(e, v), nil
	}, nil)

	assert.Equal(t, []int{1, 1, 3, 1}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.NoError(t, err)
}

func Test_Slice_AppendMapResolvOrderr(t *testing.T) {
	kvExtractor := func(val string) (bool, int, error) {
		i, err := strconv.Atoi(val)
		return even(i), i, err
	}
	resolver := func(exists bool, k bool, vr []int, v int) ([]int, error) { return resolv.Slice(exists, k, vr, v), nil }

	order, groups, err := slice.AppendMapResolvOrderr(slice.Of("2", "1", "1", "2", "4", "3", "_1"), kvExtractor, resolver, nil, nil)

	assert.Equal(t, []int{1, 1, 3}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.Equal(t, []bool{true, false}, order)
	assert.EqualError(t, err, "strconv.Atoi: parsing \"_1\": invalid syntax")

	order, groups, err = slice.AppendMapResolvOrderr(slice.Of("2", "1", "1", "2", "4", "3", "1"), kvExtractor, resolver, nil, nil)
	assert.Equal(t, []int{1, 1, 3, 1}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.Equal(t, []bool{true, false}, order)
	assert.NoError(t, err)

	resolver = func(exists bool, k bool, vr []int, v int) ([]int, error) {
		return resolv.Slice(exists, k, vr, v), op.IfElse(v == 5, errors.New("abort"), nil)
	}
	order, groups, err = slice.AppendMapResolvOrderr(slice.Of("2", "1", "1", "2", "4", "3", "5"), kvExtractor, resolver, nil, nil)
	assert.Equal(t, []int{1, 1, 3}, groups[false])
	assert.Equal(t, []int{2, 2, 4}, groups[true])
	assert.Equal(t, []bool{true, false}, order)
	assert.EqualError(t, err, "abort")

}

func Test_Slice_Filled(t *testing.T) {
	assert.Nil(t, slice.Filled[[]int](nil, nil))
	assert.Equal(t, slice.Of(1, 2, 3), slice.Filled[[]int](nil, slice.Of(1, 2, 3)))
	assert.Equal(t, slice.Of(1, 2, 3), slice.Filled(slice.Of(1, 2, 3), slice.Of(1, 2)))
}

func Test_Head(t *testing.T) {
	result, ok := slice.Head([]int{1, 3, 5, 7, 9, 11})

	assert.True(t, ok)
	assert.Equal(t, 1, result)

	_, ok = slice.Head[[]int](nil)
	assert.False(t, ok)
}

func Test_Tail(t *testing.T) {
	result, ok := slice.Tail([]int{1, 3, 5, 7, 9, 11})

	assert.True(t, ok)
	assert.Equal(t, 11, result)

	_, ok = slice.Tail[[]int](nil)
	assert.False(t, ok)
}

func Test_Contains(t *testing.T) {
	s := []int{1, 3, 5, 7, 9, 11}

	assert.True(t, slice.Contains(s, 5))
	assert.False(t, slice.Contains(s, 12))
	assert.False(t, slice.Contains(([]int)(nil), 12))

}

func Test_Upcast(t *testing.T) {
	type names []string
	n := names{"Alice"}
	strings := slice.Upcast(n)
	assert.Equal(t, []string{"Alice"}, strings)
}

func Test_Downcast(t *testing.T) {
	type names []string
	s := []string{"Alice"}
	n := slice.Downcast[names](s)
	assert.Equal(t, names{"Alice"}, n)
}
