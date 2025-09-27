package test

import (
	"testing"
	
	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
	
)

func TestNew(t *testing.T) {
	t.Run("string_int_pair", func(t *testing.T) {
		kvPair := kv.New("key", 42)
		assert.Equal(t, "key", kvPair.Key())
		assert.Equal(t, 42, kvPair.Value())

		key, value := kvPair.Get()
		assert.Equal(t, "key", key)
		assert.Equal(t, 42, value)
	})

	t.Run("int_string_pair", func(t *testing.T) {
		kvPair := kv.New(123, "value")
		assert.Equal(t, 123, kvPair.Key())
		assert.Equal(t, "value", kvPair.Value())
	})

	t.Run("same_type_pair", func(t *testing.T) {
		kvPair := kv.New(10, 20)
		assert.Equal(t, 10, kvPair.Key())
		assert.Equal(t, 20, kvPair.Value())
	})

	t.Run("zero_values", func(t *testing.T) {
		kvPair := kv.New("", 0)
		assert.Equal(t, "", kvPair.Key())
		assert.Equal(t, 0, kvPair.Value())

		kvZero := kv.New(0, "")
		assert.Equal(t, 0, kvZero.Key())
		assert.Equal(t, "", kvZero.Value())
	})

	t.Run("nil_pointer_values", func(t *testing.T) {
		var ptr *int
		kvPair := kv.New("key", ptr)
		assert.Equal(t, "key", kvPair.Key())
		assert.Nil(t, kvPair.Value())

		kvNilKey := kv.New(ptr, "value")
		assert.Nil(t, kvNilKey.Key())
		assert.Equal(t, "value", kvNilKey.Value())
	})

	t.Run("struct_types", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		p1 := person{name: "Alice", age: 30}
		p2 := person{name: "Bob", age: 25}

		kvPair := kv.New(p1, p2)
		assert.Equal(t, p1, kvPair.Key())
		assert.Equal(t, p2, kvPair.Value())
	})

	t.Run("slice_types", func(t *testing.T) {
		keys := []string{"a", "b", "c"}
		values := []int{1, 2, 3}

		kvPair := kv.New(keys, values)
		assert.Equal(t, keys, kvPair.Key())
		assert.Equal(t, values, kvPair.Value())
	})

	t.Run("map_types", func(t *testing.T) {
		keyMap := map[string]int{"x": 1, "y": 2}
		valueMap := map[int]string{1: "one", 2: "two"}

		kvPair := kv.New(keyMap, valueMap)
		assert.Equal(t, keyMap, kvPair.Key())
		assert.Equal(t, valueMap, kvPair.Value())
	})

	t.Run("all_numeric_types", func(t *testing.T) {
		// Test all integer types
		kvInt8 := kv.New(int8(-128), int8(127))
		assert.Equal(t, int8(-128), kvInt8.Key())
		assert.Equal(t, int8(127), kvInt8.Value())

		kvInt16 := kv.New(int16(-32768), int16(32767))
		assert.Equal(t, int16(-32768), kvInt16.Key())
		assert.Equal(t, int16(32767), kvInt16.Value())

		kvInt32 := kv.New(int32(-2147483648), int32(2147483647))
		assert.Equal(t, int32(-2147483648), kvInt32.Key())
		assert.Equal(t, int32(2147483647), kvInt32.Value())

		kvInt64 := kv.New(int64(-9223372036854775808), int64(9223372036854775807))
		assert.Equal(t, int64(-9223372036854775808), kvInt64.Key())
		assert.Equal(t, int64(9223372036854775807), kvInt64.Value())

		// Test unsigned types
		kvUint8 := kv.New(uint8(0), uint8(255))
		assert.Equal(t, uint8(0), kvUint8.Key())
		assert.Equal(t, uint8(255), kvUint8.Value())

		kvUint16 := kv.New(uint16(0), uint16(65535))
		assert.Equal(t, uint16(0), kvUint16.Key())
		assert.Equal(t, uint16(65535), kvUint16.Value())

		kvUint32 := kv.New(uint32(0), uint32(4294967295))
		assert.Equal(t, uint32(0), kvUint32.Key())
		assert.Equal(t, uint32(4294967295), kvUint32.Value())

		kvUint64 := kv.New(uint64(0), uint64(18446744073709551615))
		assert.Equal(t, uint64(0), kvUint64.Key())
		assert.Equal(t, uint64(18446744073709551615), kvUint64.Value())

		// Test floating point types
		kvFloat32 := kv.New(float32(-3.14), float32(3.14))
		assert.Equal(t, float32(-3.14), kvFloat32.Key())
		assert.Equal(t, float32(3.14), kvFloat32.Value())

		kvFloat64 := kv.New(-3.141592653589793, 3.141592653589793)
		assert.Equal(t, -3.141592653589793, kvFloat64.Key())
		assert.Equal(t, 3.141592653589793, kvFloat64.Value())

		// Test character types
		kvRune := kv.New('α', 'ω')
		assert.Equal(t, 'α', kvRune.Key())
		assert.Equal(t, 'ω', kvRune.Value())

		kvByte := kv.New(byte(0), byte(255))
		assert.Equal(t, byte(0), kvByte.Key())
		assert.Equal(t, byte(255), kvByte.Value())
	})

	t.Run("interface_compatibility", func(t *testing.T) {
		kvPair := kv.New("test", 42)

		// Test that it implements c.KV interface
		var kvInterface c.KV[string, int] = kvPair
		assert.Equal(t, "test", kvInterface.Key())
		assert.Equal(t, 42, kvInterface.Value())
	})

	t.Run("function_types", func(t *testing.T) {
		keyFunc := func() string { return "generated_key" }
		valueFunc := func(x int) int { return x * 2 }

		kvPair := kv.New(keyFunc, valueFunc)
		assert.Equal(t, "generated_key", kvPair.Key()())
		assert.Equal(t, 10, kvPair.Value()(5))
	})

	t.Run("channel_types", func(t *testing.T) {
		keyChan := make(chan string, 1)
		valueChan := make(chan int, 1)

		keyChan <- "channel_key"
		valueChan <- 123

		kvPair := kv.New(keyChan, valueChan)
		assert.Equal(t, "channel_key", <-kvPair.Key())
		assert.Equal(t, 123, <-kvPair.Value())
	})

	t.Run("unicode_strings", func(t *testing.T) {
		// Test various Unicode characters
		kvEmoji := kv.New("🔑", "🎯")
		assert.Equal(t, "🔑", kvEmoji.Key())
		assert.Equal(t, "🎯", kvEmoji.Value())

		kvCyrillic := kv.New("ключ", "значение")
		assert.Equal(t, "ключ", kvCyrillic.Key())
		assert.Equal(t, "значение", kvCyrillic.Value())

		kvCJK := kv.New("键", "值")
		assert.Equal(t, "键", kvCJK.Key())
		assert.Equal(t, "值", kvCJK.Value())

		kvArabic := kv.New("مفتاح", "قيمة")
		assert.Equal(t, "مفتاح", kvArabic.Key())
		assert.Equal(t, "قيمة", kvArabic.Value())
	})

	t.Run("extreme_values", func(t *testing.T) {
		// Test with very large strings
		largeKey := string(make([]rune, 100000))
		for i := range largeKey {
			largeKey = string(append([]rune(largeKey[:i]), 'k')) + largeKey[i+1:]
		}
		largeValue := string(make([]rune, 100000))
		for i := range largeValue {
			largeValue = string(append([]rune(largeValue[:i]), 'v')) + largeValue[i+1:]
		}

		kvLarge := kv.New(largeKey, largeValue)
		assert.Equal(t, largeKey, kvLarge.Key())
		assert.Equal(t, largeValue, kvLarge.Value())
	})

	t.Run("nested_structures", func(t *testing.T) {
		type address struct {
			street string
			city   string
		}

		type person struct {
			name    string
			age     int
			address address
		}

		addr := address{street: "123 Main St", city: "Anytown"}
		p := person{name: "John Doe", age: 30, address: addr}

		kvNested := kv.New("person", p)
		assert.Equal(t, "person", kvNested.Key())
		assert.Equal(t, p, kvNested.Value())
		assert.Equal(t, "123 Main St", kvNested.Value().address.street)
	})
}
