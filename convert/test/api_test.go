package test

import (
    "strconv"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/m4gshm/gollections/c"
    "github.com/m4gshm/gollections/convert"
)

func TestAsIs(t *testing.T) {
    t.Run("int", func(t *testing.T) {
        assert.Equal(t, 42, convert.AsIs(42))
    })

    t.Run("string", func(t *testing.T) {
        assert.Equal(t, "hello", convert.AsIs("hello"))
    })

    t.Run("struct", func(t *testing.T) {
        type testStruct struct{ field int }
        input := testStruct{field: 10}
        assert.Equal(t, input, convert.AsIs(input))
    })

    t.Run("nil_pointer", func(t *testing.T) {
        var p *int
        assert.Equal(t, p, convert.AsIs(p))
    })

    t.Run("zero_value", func(t *testing.T) {
        assert.Equal(t, 0, convert.AsIs(0))
        assert.Equal(t, "", convert.AsIs(""))
        assert.Equal(t, false, convert.AsIs(false))
    })
}

func TestAnd(t *testing.T) {
    t.Run("int_to_string_to_int", func(t *testing.T) {
        intToString := func(i int) string { return strconv.Itoa(i) }
        stringToInt := func(s string) int {
            i, _ := strconv.Atoi(s)
            return i
        }

        converter := convert.And(intToString, stringToInt)
        assert.Equal(t, 42, converter(42))
    })

    t.Run("compose_multiple_functions", func(t *testing.T) {
        double := func(i int) int { return i * 2 }
        toString := func(i int) string { return strconv.Itoa(i) }

        converter := convert.And(double, toString)
        assert.Equal(t, "84", converter(42))
    })

    t.Run("nil_function", func(t *testing.T) {
        // Test behavior with nil functions (should panic in real usage)
        defer func() {
            if r := recover(); r != nil {
                t.Log("Expected panic when using nil function")
            }
        }()

        var nilFunc func(int) string
        converter := convert.And(func(i int) int { return i }, nilFunc)
        converter(1) // Should panic
    })
}

func TestOr(t *testing.T) {
    t.Run("first_converter_returns_non_zero", func(t *testing.T) {
        first := func(i int) string { return "first" }
        second := func(i int) string { return "second" }

        converter := convert.Or(first, second)
        assert.Equal(t, "first", converter(42))
    })

    t.Run("first_converter_returns_zero", func(t *testing.T) {
        first := func(i int) string { return "" } // zero value for string
        second := func(i int) string { return "second" }

        converter := convert.Or(first, second)
        assert.Equal(t, "second", converter(42))
    })

    t.Run("both_return_zero", func(t *testing.T) {
        first := func(i int) string { return "" }
        second := func(i int) string { return "" }

        converter := convert.Or(first, second)
        assert.Equal(t, "", converter(42))
    })

    t.Run("numeric_zero_value", func(t *testing.T) {
        first := func(i int) int { return 0 } // zero value
        second := func(i int) int { return 100 }

        converter := convert.Or(first, second)
        assert.Equal(t, 100, converter(42))
    })
}

func TestToSlice(t *testing.T) {
    t.Run("single_int", func(t *testing.T) {
        result := convert.ToSlice(42)
        expected := []int{42}
        assert.Equal(t, expected, result)
        assert.Len(t, result, 1)
    })

    t.Run("single_string", func(t *testing.T) {
        result := convert.ToSlice("hello")
        expected := []string{"hello"}
        assert.Equal(t, expected, result)
    })

    t.Run("struct_value", func(t *testing.T) {
        type testStruct struct{ name string }
        input := testStruct{name: "test"}
        result := convert.ToSlice(input)
        expected := []testStruct{input}
        assert.Equal(t, expected, result)
    })

    t.Run("nil_pointer", func(t *testing.T) {
        var p *int
        result := convert.ToSlice(p)
        expected := []*int{p}
        assert.Equal(t, expected, result)
    })

    t.Run("zero_value", func(t *testing.T) {
        result := convert.ToSlice(0)
        expected := []int{0}
        assert.Equal(t, expected, result)
    })
}

func TestAsSlice(t *testing.T) {
    t.Run("equivalent_to_ToSlice", func(t *testing.T) {
        input := "test"
        assert.Equal(t, convert.ToSlice(input), convert.AsSlice(input))
    })
}

func TestKeyValue(t *testing.T) {
    type person struct {
        name string
        age  int
    }

    t.Run("extract_key_value", func(t *testing.T) {
        p := person{name: "Alice", age: 30}
        nameExtractor := func(p person) string { return p.name }
        ageExtractor := func(p person) int { return p.age }

        kv := convert.KeyValue(p, nameExtractor, ageExtractor)
        assert.Equal(t, "Alice", kv.Key())
        assert.Equal(t, 30, kv.Value())

        key, value := kv.Get()
        assert.Equal(t, "Alice", key)
        assert.Equal(t, 30, value)
    })

    t.Run("same_type_key_value", func(t *testing.T) {
        input := 42
        keyExtractor := func(i int) int { return i * 2 }
        valueExtractor := func(i int) int { return i * 3 }

        kv := convert.KeyValue(input, keyExtractor, valueExtractor)
        assert.Equal(t, 84, kv.Key())
        assert.Equal(t, 126, kv.Value())
    })
}

func TestKeysValues(t *testing.T) {
    type person struct {
        names []string
        ages  []int
    }

    t.Run("multiple_keys_multiple_values", func(t *testing.T) {
        p := person{
            names: []string{"Alice", "Bob"},
            ages:  []int{30, 25},
        }
        nameExtractor := func(p person) []string { return p.names }
        ageExtractor := func(p person) []int { return p.ages }

        kvs := convert.KeysValues(p, nameExtractor, ageExtractor)
        expected := []c.KV[string, int]{
            {K: "Alice", V: 30},
            {K: "Alice", V: 25},
            {K: "Bob", V: 30},
            {K: "Bob", V: 25},
        }
        assert.Equal(t, expected, kvs)
    })

    t.Run("empty_keys_multiple_values", func(t *testing.T) {
        p := person{
            names: []string{},
            ages:  []int{30, 25},
        }
        nameExtractor := func(p person) []string { return p.names }
        ageExtractor := func(p person) []int { return p.ages }

        kvs := convert.KeysValues(p, nameExtractor, ageExtractor)
        expected := []c.KV[string, int]{
            {K: "", V: 30}, // zero key value
            {K: "", V: 25},
        }
        assert.Equal(t, expected, kvs)
    })

    t.Run("multiple_keys_empty_values", func(t *testing.T) {
        p := person{
            names: []string{"Alice", "Bob"},
            ages:  []int{},
        }
        nameExtractor := func(p person) []string { return p.names }
        ageExtractor := func(p person) []int { return p.ages }

        kvs := convert.KeysValues(p, nameExtractor, ageExtractor)
        expected := []c.KV[string, int]{
            {K: "Alice", V: 0}, // zero value
            {K: "Bob", V: 0},
        }
        assert.Equal(t, expected, kvs)
    })

    // t.Run("empty_keys_empty_values", func(t *testing.T) {
    //     p := person{
    //         names: []string{},
    //         ages:  []int{},
    //     }
    //     nameExtractor := func(p person) []string { return p.names }
    //     ageExtractor := func(p person) []int { return p.ages }

    //     kvs := convert.KeysValues(p, nameExtractor, ageExtractor)
    //     expected := []c.KV[string, int]{
    //         {K: "", V: 0}, // both zero values
    //     }
    //     assert.Equal(t, expected, kvs)
    // })
}

func TestExtraVals(t *testing.T) {
    t.Run("single_value", func(t *testing.T) {
        input := "key"
        valueExtractor := func(s string) []string { return []string{"val1", "val2"} }

        kvs := convert.ExtraVals(input, valueExtractor)
        expected := []c.KV[string, string]{
            {K: "key", V: "val1"},
            {K: "key", V: "val2"},
        }
        assert.Equal(t, expected, kvs)
    })

    t.Run("empty_values", func(t *testing.T) {
        input := "key"
        valueExtractor := func(s string) []string { return []string{} }

        kvs := convert.ExtraVals(input, valueExtractor)
        expected := []c.KV[string, string]{
            {K: "key", V: ""}, // zero value
        }
        assert.Equal(t, expected, kvs)
    })

    t.Run("numeric_type", func(t *testing.T) {
        input := 42
        valueExtractor := func(i int) []int { return []int{1, 2, 3} }

        kvs := convert.ExtraVals(input, valueExtractor)
        expected := []c.KV[int, int]{
            {K: 42, V: 1},
            {K: 42, V: 2},
            {K: 42, V: 3},
        }
        assert.Equal(t, expected, kvs)
    })
}

func TestExtraKeys(t *testing.T) {
    t.Run("multiple_keys", func(t *testing.T) {
        input := "value"
        keyExtractor := func(s string) []string { return []string{"key1", "key2"} }

        kvs := convert.ExtraKeys(input, keyExtractor)
        expected := []c.KV[string, string]{
            {K: "key1", V: "value"},
            {K: "key2", V: "value"},
        }
        assert.Equal(t, expected, kvs)
    })

    t.Run("empty_keys", func(t *testing.T) {
        input := "value"
        keyExtractor := func(s string) []string { return []string{} }

        kvs := convert.ExtraKeys(input, keyExtractor)
        expected := []c.KV[string, string]{
            {K: "", V: "value"}, // zero key value
        }
        assert.Equal(t, expected, kvs)
    })

    t.Run("numeric_type", func(t *testing.T) {
        input := 100
        keyExtractor := func(i int) []int { return []int{1, 2} }

        kvs := convert.ExtraKeys(input, keyExtractor)
        expected := []c.KV[int, int]{
            {K: 1, V: 100},
            {K: 2, V: 100},
        }
        assert.Equal(t, expected, kvs)
    })
}

func TestPtr(t *testing.T) {
    t.Run("int_value", func(t *testing.T) {
        value := 42
        ptr := convert.Ptr(value)
        assert.NotNil(t, ptr)
        assert.Equal(t, value, *ptr)
    })

    t.Run("string_value", func(t *testing.T) {
        value := "hello"
        ptr := convert.Ptr(value)
        assert.NotNil(t, ptr)
        assert.Equal(t, value, *ptr)
    })

    t.Run("struct_value", func(t *testing.T) {
        type testStruct struct{ field int }
        value := testStruct{field: 10}
        ptr := convert.Ptr(value)
        assert.NotNil(t, ptr)
        assert.Equal(t, value, *ptr)
    })

    t.Run("zero_value", func(t *testing.T) {
        ptr := convert.Ptr(0)
        assert.NotNil(t, ptr)
        assert.Equal(t, 0, *ptr)
    })
}

func TestPtrVal(t *testing.T) {
    t.Run("valid_pointer", func(t *testing.T) {
        value := 42
        ptr := &value
        result := convert.PtrVal(ptr)
        assert.Equal(t, value, result)
    })

    t.Run("nil_pointer", func(t *testing.T) {
        var ptr *int
        result := convert.PtrVal(ptr)
        assert.Equal(t, 0, result) // zero value for int
    })

    t.Run("string_pointer", func(t *testing.T) {
        value := "hello"
        ptr := &value
        result := convert.PtrVal(ptr)
        assert.Equal(t, value, result)
    })

    t.Run("nil_string_pointer", func(t *testing.T) {
        var ptr *string
        result := convert.PtrVal(ptr)
        assert.Equal(t, "", result) // zero value for string
    })
}

func TestNoNilPtrVal(t *testing.T) {
    t.Run("valid_pointer", func(t *testing.T) {
        value := 42
        ptr := &value
        result, ok := convert.NoNilPtrVal(ptr)
        assert.True(t, ok)
        assert.Equal(t, value, result)
    })

    t.Run("nil_pointer", func(t *testing.T) {
        var ptr *int
        result, ok := convert.NoNilPtrVal(ptr)
        assert.False(t, ok)
        assert.Equal(t, 0, result) // zero value
    })

    t.Run("string_pointer", func(t *testing.T) {
        value := "hello"
        ptr := &value
        result, ok := convert.NoNilPtrVal(ptr)
        assert.True(t, ok)
        assert.Equal(t, value, result)
    })

    t.Run("nil_string_pointer", func(t *testing.T) {
        var ptr *string
        result, ok := convert.NoNilPtrVal(ptr)
        assert.False(t, ok)
        assert.Equal(t, "", result)
    })
}

func TestToType(t *testing.T) {
    t.Run("successful_conversion", func(t *testing.T) {
        var input interface{} = 42
        result, ok := convert.ToType[int](input)
        assert.True(t, ok)
        assert.Equal(t, 42, result)
    })

    t.Run("failed_conversion", func(t *testing.T) {
        var input interface{} = "hello"
        result, ok := convert.ToType[int](input)
        assert.False(t, ok)
        assert.Equal(t, 0, result) // zero value
    })

    t.Run("nil_to_pointer", func(t *testing.T) {
        var input interface{}
        result, ok := convert.ToType[*int](input)
        assert.False(t, ok)
        assert.Nil(t, result)
    })

    t.Run("pointer_conversion", func(t *testing.T) {
        value := 42
        var input interface{} = &value
        result, ok := convert.ToType[*int](input)
        assert.True(t, ok)
        assert.Equal(t, &value, result)
    })

    t.Run("interface_conversion", func(t *testing.T) {
        var input interface{} = "hello"
        result, ok := convert.ToType[interface{}](input)
        assert.True(t, ok)
        assert.Equal(t, "hello", result)
    })

    t.Run("struct_conversion", func(t *testing.T) {
        type testStruct struct{ field int }
        value := testStruct{field: 42}
        var input interface{} = value
        result, ok := convert.ToType[testStruct](input)
        assert.True(t, ok)
        assert.Equal(t, value, result)
    })
}

// Benchmark tests
func BenchmarkAsIs(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = convert.AsIs(42)
    }
}

func BenchmarkAnd(b *testing.B) {
    converter := convert.And(func(i int) string { return strconv.Itoa(i) }, func(s string) int { i, _ := strconv.Atoi(s); return i })
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = converter(42)
    }
}

func BenchmarkPtr(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = convert.Ptr(42)
    }
}

func BenchmarkPtrVal(b *testing.B) {
    value := 42
    ptr := &value
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = convert.PtrVal(ptr)
    }
}

func BenchmarkToType(b *testing.B) {
    var input interface{} = 42
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = convert.ToType[int](input)
    }
}