package optional

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Test constructors
func TestSome(t *testing.T) {
	t.Run("Create Some with int", func(t *testing.T) {
		opt := Some(42)
		if !opt.IsPresent() {
			t.Error("Some should be present")
		}
		if opt.Get() != 42 {
			t.Errorf("Expected 42, got %v", opt.Get())
		}
	})

	t.Run("Create Some with string", func(t *testing.T) {
		opt := Some("hello")
		if !opt.IsPresent() {
			t.Error("Some should be present")
		}
		if opt.Get() != "hello" {
			t.Errorf("Expected 'hello', got %v", opt.Get())
		}
	})

	t.Run("Create Some with nil pointer", func(t *testing.T) {
		var ptr *string
		opt := Some(ptr)
		if !opt.IsPresent() {
			t.Error("Some should be present even with nil pointer")
		}
		if opt.Get() != nil {
			t.Error("Expected nil pointer")
		}
	})
}

func TestNone(t *testing.T) {
	t.Run("Create None with int", func(t *testing.T) {
		opt := None[int]()
		if opt.IsPresent() {
			t.Error("None should not be present")
		}
		if !opt.IsEmpty() {
			t.Error("None should be empty")
		}
	})

	t.Run("Create None with string", func(t *testing.T) {
		opt := None[string]()
		if opt.IsPresent() {
			t.Error("None should not be present")
		}
		if !opt.IsEmpty() {
			t.Error("None should be empty")
		}
	})
}

func TestFromPointer(t *testing.T) {
	t.Run("FromPointer with non-nil pointer", func(t *testing.T) {
		value := 42
		opt := FromPointer(&value)
		if !opt.IsPresent() {
			t.Error("FromPointer with non-nil should be present")
		}
		if opt.Get() != 42 {
			t.Errorf("Expected 42, got %v", opt.Get())
		}
	})

	t.Run("FromPointer with nil pointer", func(t *testing.T) {
		var ptr *int
		opt := FromPointer(ptr)
		if opt.IsPresent() {
			t.Error("FromPointer with nil should not be present")
		}
	})
}

func TestToPointer(t *testing.T) {
	t.Run("ToPointer with present value", func(t *testing.T) {
		opt := Some(42)
		ptr := opt.ToPointer()
		if ptr == nil {
			t.Error("ToPointer should return non-nil pointer for present value")
		}
		if *ptr != 42 {
			t.Errorf("Expected 42, got %v", *ptr)
		}
	})

	t.Run("ToPointer with empty value", func(t *testing.T) {
		opt := None[int]()
		ptr := opt.ToPointer()
		if ptr != nil {
			t.Error("ToPointer should return nil for empty value")
		}
	})
}

// Test status methods
func TestIsPresent(t *testing.T) {
	t.Run("IsPresent with Some", func(t *testing.T) {
		opt := Some(42)
		if !opt.IsPresent() {
			t.Error("Some should be present")
		}
	})

	t.Run("IsPresent with None", func(t *testing.T) {
		opt := None[int]()
		if opt.IsPresent() {
			t.Error("None should not be present")
		}
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty with Some", func(t *testing.T) {
		opt := Some(42)
		if opt.IsEmpty() {
			t.Error("Some should not be empty")
		}
	})

	t.Run("IsEmpty with None", func(t *testing.T) {
		opt := None[int]()
		if !opt.IsEmpty() {
			t.Error("None should be empty")
		}
	})
}

// Test value retrieval
func TestGet(t *testing.T) {
	t.Run("Get from Some", func(t *testing.T) {
		opt := Some(42)
		value := opt.Get()
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("Get from None should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Get from None should panic")
			}
		}()
		opt := None[int]()
		opt.Get()
	})
}

func TestEquals(t *testing.T) {
	t.Run("Equal Some values", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := Some(42)
		if !opt1.Equals(opt2) {
			t.Error("Equal Some values should be equal")
		}
	})

	t.Run("Different Some values", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := Some(24)
		if opt1.Equals(opt2) {
			t.Error("Different Some values should not be equal")
		}
	})

	t.Run("Equal None values", func(t *testing.T) {
		opt1 := None[int]()
		opt2 := None[int]()
		if !opt1.Equals(opt2) {
			t.Error("None values should be equal")
		}
	})

	t.Run("Some vs None", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := None[int]()
		if opt1.Equals(opt2) {
			t.Error("Some and None should not be equal")
		}
	})
}

func TestOr(t *testing.T) {
	t.Run("Some.Or(Some)", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := Some(24)
		result := opt1.Or(opt2)
		if !result.IsPresent() || result.Get() != 42 {
			t.Error("Some.Or should return first Some")
		}
	})

	t.Run("Some.Or(None)", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := None[int]()
		result := opt1.Or(opt2)
		if !result.IsPresent() || result.Get() != 42 {
			t.Error("Some.Or should return first Some")
		}
	})

	t.Run("None.Or(Some)", func(t *testing.T) {
		opt1 := None[int]()
		opt2 := Some(42)
		result := opt1.Or(opt2)
		if !result.IsPresent() || result.Get() != 42 {
			t.Error("None.Or should return second Optional")
		}
	})

	t.Run("None.Or(None)", func(t *testing.T) {
		opt1 := None[int]()
		opt2 := None[int]()
		result := opt1.Or(opt2)
		if result.IsPresent() {
			t.Error("None.Or(None) should return None")
		}
	})
}

func TestOrElsePanic(t *testing.T) {
	t.Run("OrElsePanic with Some", func(t *testing.T) {
		opt := Some(42)
		value := opt.OrElsePanic("should not panic")
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("OrElsePanic with None should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("OrElsePanic with None should panic")
			} else if r != "custom panic message" {
				t.Errorf("Expected custom panic message, got %v", r)
			}
		}()
		opt := None[int]()
		opt.OrElsePanic("custom panic message")
	})
}

func TestOrElse(t *testing.T) {
	t.Run("OrElse with Some", func(t *testing.T) {
		opt := Some(42)
		value := opt.OrElse(0)
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
	})

	t.Run("OrElse with None", func(t *testing.T) {
		opt := None[int]()
		value := opt.OrElse(100)
		if value != 100 {
			t.Errorf("Expected 100, got %v", value)
		}
	})
}

func TestOrElseGet(t *testing.T) {
	t.Run("OrElseGet with Some", func(t *testing.T) {
		opt := Some(42)
		called := false
		value := opt.OrElseGet(func() int {
			called = true
			return 100
		})
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}
		if called {
			t.Error("Supplier should not be called for Some")
		}
	})

	t.Run("OrElseGet with None", func(t *testing.T) {
		opt := None[int]()
		called := false
		value := opt.OrElseGet(func() int {
			called = true
			return 100
		})
		if value != 100 {
			t.Errorf("Expected 100, got %v", value)
		}
		if !called {
			t.Error("Supplier should be called for None")
		}
	})
}

func TestIfPresent(t *testing.T) {
	t.Run("IfPresent with Some", func(t *testing.T) {
		opt := Some(42)
		called := false
		opt.IfPresent(func(value int) {
			called = true
			if value != 42 {
				t.Errorf("Expected 42, got %v", value)
			}
		})
		if !called {
			t.Error("Consumer should be called for Some")
		}
	})

	t.Run("IfPresent with None", func(t *testing.T) {
		opt := None[int]()
		called := false
		opt.IfPresent(func(value int) {
			called = true
		})
		if called {
			t.Error("Consumer should not be called for None")
		}
	})
}

func TestIfPresentOrElse(t *testing.T) {
	t.Run("IfPresentOrElse with Some", func(t *testing.T) {
		opt := Some(42)
		consumerCalled := false
		runnableCalled := false

		opt.IfPresentOrElse(
			func(value int) {
				consumerCalled = true
				if value != 42 {
					t.Errorf("Expected 42, got %v", value)
				}
			},
			func() {
				runnableCalled = true
			},
		)

		if !consumerCalled {
			t.Error("Consumer should be called for Some")
		}
		if runnableCalled {
			t.Error("Runnable should not be called for Some")
		}
	})

	t.Run("IfPresentOrElse with None", func(t *testing.T) {
		opt := None[int]()
		consumerCalled := false
		runnableCalled := false

		opt.IfPresentOrElse(
			func(value int) {
				consumerCalled = true
			},
			func() {
				runnableCalled = true
			},
		)

		if consumerCalled {
			t.Error("Consumer should not be called for None")
		}
		if !runnableCalled {
			t.Error("Runnable should be called for None")
		}
	})
}

func TestZip(t *testing.T) {
	t.Run("Zip two Some values", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := Some("hello")
		result := Zip(opt1, opt2, func(a int, b string) string {
			return fmt.Sprintf("%d-%s", a, b)
		})

		if !result.IsPresent() {
			t.Error("Zip of two Some should be present")
		}
		if result.Get() != "42-hello" {
			t.Errorf("Expected '42-hello', got %v", result.Get())
		}
	})

	t.Run("Zip Some and None", func(t *testing.T) {
		opt1 := Some(42)
		opt2 := None[string]()
		result := Zip(opt1, opt2, func(a int, b string) string {
			return fmt.Sprintf("%d-%s", a, b)
		})

		if result.IsPresent() {
			t.Error("Zip of Some and None should be None")
		}
	})

	t.Run("Zip None and Some", func(t *testing.T) {
		opt1 := None[int]()
		opt2 := Some("hello")
		result := Zip(opt1, opt2, func(a int, b string) string {
			return fmt.Sprintf("%d-%s", a, b)
		})

		if result.IsPresent() {
			t.Error("Zip of None and Some should be None")
		}
	})

	t.Run("Zip two None values", func(t *testing.T) {
		opt1 := None[int]()
		opt2 := None[string]()
		result := Zip(opt1, opt2, func(a int, b string) string {
			return fmt.Sprintf("%d-%s", a, b)
		})

		if result.IsPresent() {
			t.Error("Zip of two None should be None")
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("Map with Some", func(t *testing.T) {
		opt := Some(42)
		result := Map(opt, func(x int) string {
			return fmt.Sprintf("value-%d", x)
		})

		if !result.IsPresent() {
			t.Error("Map of Some should be present")
		}
		if result.Get() != "value-42" {
			t.Errorf("Expected 'value-42', got %v", result.Get())
		}
	})

	t.Run("Map with None", func(t *testing.T) {
		opt := None[int]()
		mapperCalled := false
		result := Map(opt, func(x int) string {
			mapperCalled = true
			return fmt.Sprintf("value-%d", x)
		})

		if result.IsPresent() {
			t.Error("Map of None should be None")
		}
		if mapperCalled {
			t.Error("Mapper should not be called for None")
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("FlatMap with Some returning Some", func(t *testing.T) {
		opt := Some(42)
		result := FlatMap(opt, func(x int) Optional[string] {
			if x > 0 {
				return Some(fmt.Sprintf("positive-%d", x))
			}
			return None[string]()
		})

		if !result.IsPresent() {
			t.Error("FlatMap of Some returning Some should be present")
		}
		if result.Get() != "positive-42" {
			t.Errorf("Expected 'positive-42', got %v", result.Get())
		}
	})

	t.Run("FlatMap with Some returning None", func(t *testing.T) {
		opt := Some(-42)
		result := FlatMap(opt, func(x int) Optional[string] {
			if x > 0 {
				return Some(fmt.Sprintf("positive-%d", x))
			}
			return None[string]()
		})

		if result.IsPresent() {
			t.Error("FlatMap of Some returning None should be None")
		}
	})

	t.Run("FlatMap with None", func(t *testing.T) {
		opt := None[int]()
		mapperCalled := false
		result := FlatMap(opt, func(x int) Optional[string] {
			mapperCalled = true
			return Some(fmt.Sprintf("value-%d", x))
		})

		if result.IsPresent() {
			t.Error("FlatMap of None should be None")
		}
		if mapperCalled {
			t.Error("Mapper should not be called for None")
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("Filter with Some matching predicate", func(t *testing.T) {
		opt := Some(42)
		result := opt.Filter(func(x int) bool {
			return x > 0
		})

		if !result.IsPresent() {
			t.Error("Filter with matching predicate should be present")
		}
		if result.Get() != 42 {
			t.Errorf("Expected 42, got %v", result.Get())
		}
	})

	t.Run("Filter with Some not matching predicate", func(t *testing.T) {
		opt := Some(-42)
		result := opt.Filter(func(x int) bool {
			return x > 0
		})

		if result.IsPresent() {
			t.Error("Filter with non-matching predicate should be None")
		}
	})

	t.Run("Filter with None", func(t *testing.T) {
		opt := None[int]()
		predicateCalled := false
		result := opt.Filter(func(x int) bool {
			predicateCalled = true
			return true
		})

		if result.IsPresent() {
			t.Error("Filter of None should be None")
		}
		if predicateCalled {
			t.Error("Predicate should not be called for None")
		}
	})
}

func TestString(t *testing.T) {
	t.Run("String representation of Some", func(t *testing.T) {
		opt := Some(42)
		str := opt.String()
		if str != "Some(42)" {
			t.Errorf("Expected 'Some(42)', got %s", str)
		}
	})

	t.Run("String representation of None", func(t *testing.T) {
		opt := None[int]()
		str := opt.String()
		if str != "None" {
			t.Errorf("Expected 'None', got %s", str)
		}
	})

	t.Run("String representation of Some with string", func(t *testing.T) {
		opt := Some("hello")
		str := opt.String()
		if str != "Some(hello)" {
			t.Errorf("Expected 'Some(hello)', got %s", str)
		}
	})
}

// Test JSON serialization
func TestMarshalJSON(t *testing.T) {
	t.Run("Marshal Some", func(t *testing.T) {
		opt := Some(42)
		data, err := json.Marshal(&opt)
		if err != nil {
			t.Errorf("Marshal error: %v", err)
		}
		if string(data) != "42" {
			t.Errorf("Expected '42', got %s", string(data))
		}
	})

	t.Run("Marshal None", func(t *testing.T) {
		opt := None[int]()
		data, err := json.Marshal(&opt)
		if err != nil {
			t.Errorf("Marshal error: %v", err)
		}
		if string(data) != "null" {
			t.Errorf("Expected 'null', got %s", string(data))
		}
	})

	t.Run("Marshal Some with string", func(t *testing.T) {
		opt := Some("hello")
		data, err := json.Marshal(&opt)
		if err != nil {
			t.Errorf("Marshal error: %v", err)
		}
		if string(data) != `"hello"` {
			t.Errorf("Expected '\"hello\"', got %s", string(data))
		}
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("Unmarshal non-null value", func(t *testing.T) {
		var opt Optional[int]
		err := json.Unmarshal([]byte("42"), &opt)
		if err != nil {
			t.Errorf("Unmarshal error: %v", err)
		}
		if !opt.IsPresent() {
			t.Error("Unmarshaled value should be present")
		}
		if opt.Get() != 42 {
			t.Errorf("Expected 42, got %v", opt.Get())
		}
	})

	t.Run("Unmarshal null value", func(t *testing.T) {
		var opt Optional[int]
		err := json.Unmarshal([]byte("null"), &opt)
		if err != nil {
			t.Errorf("Unmarshal error: %v", err)
		}
		if opt.IsPresent() {
			t.Error("Unmarshaled null should not be present")
		}
	})

	t.Run("Unmarshal string value", func(t *testing.T) {
		var opt Optional[string]
		err := json.Unmarshal([]byte(`"hello"`), &opt)
		if err != nil {
			t.Errorf("Unmarshal error: %v", err)
		}
		if !opt.IsPresent() {
			t.Error("Unmarshaled value should be present")
		}
		if opt.Get() != "hello" {
			t.Errorf("Expected 'hello', got %v", opt.Get())
		}
	})

	t.Run("Unmarshal invalid JSON", func(t *testing.T) {
		var opt Optional[int]
		err := json.Unmarshal([]byte("invalid"), &opt)
		if err == nil {
			t.Error("Should return error for invalid JSON")
		}
	})
}

// Integration tests
func TestIntegration(t *testing.T) {
	t.Run("Chain operations", func(t *testing.T) {
		result := Some(42)
		mapped := Map(result, func(x int) int { return x * 2 })
		filtered := mapped.Filter(func(x int) bool { return x > 50 })
		final := filtered.OrElse(0)

		if final != 84 {
			t.Errorf("Expected 84, got %v", final)
		}
	})

	t.Run("JSON roundtrip", func(t *testing.T) {
		original := Some("test value")

		// Marshal
		data, err := json.Marshal(&original)
		if err != nil {
			t.Errorf("Marshal error: %v", err)
		}

		// Unmarshal
		var restored Optional[string]
		err = json.Unmarshal(data, &restored)
		if err != nil {
			t.Errorf("Unmarshal error: %v", err)
		}

		// Compare
		if !original.Equals(restored) {
			t.Error("JSON roundtrip failed")
		}
	})
}

// Existing benchmarks
func BenchmarkSome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Some(42)
	}
}

func BenchmarkNone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		None[int]()
	}
}

func BenchmarkMap(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(opt, func(x int) int { return x * 2 })
	}
}

func BenchmarkFilter(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opt.Filter(func(x int) bool { return x > 0 })
	}
}

// Additional detailed benchmarks
func BenchmarkGet(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.Get()
	}
}

func BenchmarkOrElse(b *testing.B) {
	b.Run("Some", func(b *testing.B) {
		opt := Some(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = opt.OrElse(0)
		}
	})

	b.Run("None", func(b *testing.B) {
		opt := None[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = opt.OrElse(100)
		}
	})
}

func BenchmarkIsPresent(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.IsPresent()
	}
}

func BenchmarkFlatMap(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FlatMap(opt, func(x int) Optional[int] {
			if x > 0 {
				return Some(x * 2)
			}
			return None[int]()
		})
	}
}

func BenchmarkChainOperations(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Some(42)
		result = Map(result, func(x int) int { return x * 2 })
		result = result.Filter(func(x int) bool { return x > 50 })
		_ = result.OrElse(0)
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&opt)
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	data := []byte("42")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var opt Optional[int]
		_ = json.Unmarshal(data, &opt)
	}
}

// So sánh với pointer operations
func BenchmarkPointerComparison(b *testing.B) {
	b.Run("Optional", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			opt := Some(42)
			if opt.IsPresent() {
				_ = opt.Get()
			}
		}
	})

	b.Run("Pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := 42
			ptr := &value
			if ptr != nil {
				_ = *ptr
			}
		}
	})
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("Some", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Some(42)
		}
	})

	b.Run("Map", func(b *testing.B) {
		opt := Some(42)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Map(opt, func(x int) int { return x * 2 })
		}
	})
}

// String performance
func BenchmarkString(b *testing.B) {
	b.Run("Some", func(b *testing.B) {
		opt := Some(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = opt.String()
		}
	})

	b.Run("None", func(b *testing.B) {
		opt := None[int]()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = opt.String()
		}
	})
}
