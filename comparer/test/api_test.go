package test

import (
    "testing"

    "github.com/stretchr/testify/assert"
	"github.com/m4gshm/gollections/comparer"
)

type person struct {
    name string
    age  int
}

func TestOf(t *testing.T) {
    t.Run("compare_by_age", func(t *testing.T) {
        ageExtractor := func(p person) int { return p.age }
        comparer := comparer.Of(ageExtractor)

        p1 := person{name: "Alice", age: 25}
        p2 := person{name: "Bob", age: 30}
        p3 := person{name: "Charlie", age: 25}

        // p1.age < p2.age
        assert.Equal(t, -1, comparer(p1, p2))
        // p2.age > p1.age
        assert.Equal(t, 1, comparer(p2, p1))
        // p1.age == p3.age
        assert.Equal(t, 0, comparer(p1, p3))
    })

    t.Run("compare_by_name", func(t *testing.T) {
        nameExtractor := func(p person) string { return p.name }
        comparer := comparer.Of(nameExtractor)

        p1 := person{name: "Alice", age: 25}
        p2 := person{name: "Bob", age: 30}
        p3 := person{name: "Alice", age: 35}

        // "Alice" < "Bob"
        assert.Equal(t, -1, comparer(p1, p2))
        // "Bob" > "Alice"
        assert.Equal(t, 1, comparer(p2, p1))
        // "Alice" == "Alice"
        assert.Equal(t, 0, comparer(p1, p3))
    })

    t.Run("compare_integers", func(t *testing.T) {
        identity := func(i int) int { return i }
        comparer := comparer.Of(identity)

        assert.Equal(t, -1, comparer(1, 2))
        assert.Equal(t, 1, comparer(2, 1))
        assert.Equal(t, 0, comparer(1, 1))
    })

    t.Run("compare_floats", func(t *testing.T) {
        identity := func(f float64) float64 { return f }
        comparer := comparer.Of(identity)

        assert.Equal(t, -1, comparer(1.5, 2.5))
        assert.Equal(t, 1, comparer(2.5, 1.5))
        assert.Equal(t, 0, comparer(1.5, 1.5))
    })

    t.Run("compare_by_computed_value", func(t *testing.T) {
        // Compare by string length
        lengthExtractor := func(s string) int { return len(s) }
        comparer := comparer.Of(lengthExtractor)

        assert.Equal(t, -1, comparer("hi", "hello"))    // 2 < 5
        assert.Equal(t, 1, comparer("hello", "hi"))     // 5 > 2
        assert.Equal(t, 0, comparer("hello", "world"))  // 5 == 5
    })

    t.Run("compare_negative_numbers", func(t *testing.T) {
        identity := func(i int) int { return i }
        comparer := comparer.Of(identity)

        assert.Equal(t, 1, comparer(-1, -2))   // -1 > -2
        assert.Equal(t, -1, comparer(-2, -1))  // -2 < -1
        assert.Equal(t, 0, comparer(-1, -1))   // -1 == -1
    })
}

func TestReverse(t *testing.T) {
    t.Run("reverse_compare_by_age", func(t *testing.T) {
        ageExtractor := func(p person) int { return p.age }
        comparer := comparer.Reverse(ageExtractor)

        p1 := person{name: "Alice", age: 25}
        p2 := person{name: "Bob", age: 30}
        p3 := person{name: "Charlie", age: 25}

        // Reverse: p1.age < p2.age becomes p1 > p2
        assert.Equal(t, 1, comparer(p1, p2))
        // Reverse: p2.age > p1.age becomes p2 < p1
        assert.Equal(t, -1, comparer(p2, p1))
        // p1.age == p3.age stays equal
        assert.Equal(t, 0, comparer(p1, p3))
    })

    t.Run("reverse_compare_by_name", func(t *testing.T) {
        nameExtractor := func(p person) string { return p.name }
        comparer := comparer.Reverse(nameExtractor)

        p1 := person{name: "Alice", age: 25}
        p2 := person{name: "Bob", age: 30}
        p3 := person{name: "Alice", age: 35}

        // Reverse: "Alice" < "Bob" becomes "Alice" > "Bob"
        assert.Equal(t, 1, comparer(p1, p2))
        // Reverse: "Bob" > "Alice" becomes "Bob" < "Alice"
        assert.Equal(t, -1, comparer(p2, p1))
        // "Alice" == "Alice" stays equal
        assert.Equal(t, 0, comparer(p1, p3))
    })

    t.Run("reverse_integers", func(t *testing.T) {
        identity := func(i int) int { return i }
        comparer := comparer.Reverse(identity)

        assert.Equal(t, 1, comparer(1, 2))   // Reverse: 1 < 2 becomes 1 > 2
        assert.Equal(t, -1, comparer(2, 1))  // Reverse: 2 > 1 becomes 2 < 1
        assert.Equal(t, 0, comparer(1, 1))   // 1 == 1 stays equal
    })

    t.Run("reverse_vs_normal_comparer", func(t *testing.T) {
        identity := func(i int) int { return i }
        normalComparer := comparer.Of(identity)
        reverseComparer := comparer.Reverse(identity)

        // Results should be opposite (except for equal values)
        assert.Equal(t, -normalComparer(1, 2), reverseComparer(1, 2))
        assert.Equal(t, -normalComparer(2, 1), reverseComparer(2, 1))
        assert.Equal(t, normalComparer(1, 1), reverseComparer(1, 1)) // Equal stays equal
    })
}

func TestComparerWithComplexTypes(t *testing.T) {
    type book struct {
        title  string
        year   int
        rating float64
    }

    t.Run("compare_books_by_year", func(t *testing.T) {
        yearExtractor := func(b book) int { return b.year }
        comparer := comparer.Of(yearExtractor)

        book1 := book{title: "Book A", year: 2020, rating: 4.5}
        book2 := book{title: "Book B", year: 2022, rating: 4.0}

        assert.Equal(t, -1, comparer(book1, book2)) // 2020 < 2022
        assert.Equal(t, 1, comparer(book2, book1))  // 2022 > 2020
    })

    t.Run("compare_books_by_rating", func(t *testing.T) {
        ratingExtractor := func(b book) float64 { return b.rating }
        comparer := comparer.Of(ratingExtractor)

        book1 := book{title: "Book A", year: 2020, rating: 4.5}
        book2 := book{title: "Book B", year: 2022, rating: 4.0}

        assert.Equal(t, 1, comparer(book1, book2))  // 4.5 > 4.0
        assert.Equal(t, -1, comparer(book2, book1)) // 4.0 < 4.5
    })
}

func TestComparerEdgeCases(t *testing.T) {
    t.Run("zero_values", func(t *testing.T) {
        identity := func(i int) int { return i }
        comparer := comparer.Of(identity)

        assert.Equal(t, -1, comparer(0, 1))  // 0 < 1
        assert.Equal(t, 1, comparer(1, 0))   // 1 > 0
        assert.Equal(t, 0, comparer(0, 0))   // 0 == 0
    })

    t.Run("empty_strings", func(t *testing.T) {
        identity := func(s string) string { return s }
        comparer := comparer.Of(identity)

        assert.Equal(t, -1, comparer("", "a"))  // "" < "a"
        assert.Equal(t, 1, comparer("a", ""))   // "a" > ""
        assert.Equal(t, 0, comparer("", ""))    // "" == ""
    })

    t.Run("unicode_strings", func(t *testing.T) {
        identity := func(s string) string { return s }
        comparer := comparer.Of(identity)

        assert.Equal(t, -1, comparer("a", "ä"))  // ASCII < Unicode
        assert.Equal(t, 1, comparer("ä", "a"))   // Unicode > ASCII
        assert.Equal(t, 0, comparer("ä", "ä"))   // Unicode == Unicode
    })
}

// Integration test with sorting behavior
func TestComparerIntegration(t *testing.T) {
    t.Run("sorting_behavior_simulation", func(t *testing.T) {
        people := []person{
            {name: "Charlie", age: 35},
            {name: "Alice", age: 25},
            {name: "Bob", age: 30},
        }

        ageComparer := comparer.Of(func(p person) int { return p.age })

        // Simulate sorting by manually comparing
        // Alice (25) should come first, then Bob (30), then Charlie (35)
        assert.Equal(t, -1, ageComparer(people[1], people[2])) // Alice < Bob
        assert.Equal(t, -1, ageComparer(people[2], people[0])) // Bob < Charlie
        assert.Equal(t, -1, ageComparer(people[1], people[0])) // Alice < Charlie
    })

    t.Run("reverse_sorting_behavior", func(t *testing.T) {
        people := []person{
            {name: "Alice", age: 25},
            {name: "Bob", age: 30},
            {name: "Charlie", age: 35},
        }

        ageReverseComparer := comparer.Reverse(func(p person) int { return p.age })

        // With reverse comparer: Charlie (35) should come first
        assert.Equal(t, 1, ageReverseComparer(people[0], people[1])) // Alice > Bob (reversed)
        assert.Equal(t, 1, ageReverseComparer(people[1], people[2])) // Bob > Charlie (reversed)
        assert.Equal(t, 1, ageReverseComparer(people[0], people[2])) // Alice > Charlie (reversed)
    })
}

// Benchmark tests
func BenchmarkOf(b *testing.B) {
    comparer := comparer.Of(func(p person) int { return p.age })
    p1 := person{name: "Alice", age: 25}
    p2 := person{name: "Bob", age: 30}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = comparer(p1, p2)
    }
}

func BenchmarkReverse(b *testing.B) {
    comparer := comparer.Reverse(func(p person) int { return p.age })
    p1 := person{name: "Alice", age: 25}
    p2 := person{name: "Bob", age: 30}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = comparer(p1, p2)
    }
}

