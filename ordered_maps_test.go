package gomaps_test

import (
	"testing"

	"github.com/dimmerz92/gomaps"
)

func TestOrderedMap(t *testing.T) {
	om := gomaps.NewOrderedMap[string, int]()

	t.Run("test push", func(t *testing.T) {
		err := om.Push("a", 1)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		err = om.Push("b", 2)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		err = om.Push("a", 1)
		if err == nil {
			t.Fatal("expected err, got nil")
		}
	})

	t.Run("test get", func(t *testing.T) {
		val, ok := om.Get("a")
		if !ok || val != 1 {
			t.Errorf("Expected (1, true), got (%d, %v)", val, ok)
		}

		val, ok = om.Get("b")
		if !ok || val != 2 {
			t.Errorf("Expected (2, true), got (%d, %v)", val, ok)
		}

		val, ok = om.Get("c")
		if ok {
			t.Errorf("Expected (0, false), got (%d, %v)", val, ok)
		}
	})

	t.Run("test set", func(t *testing.T) {
		om.Set("x", 3)
		om.Set("x", 15)

		val, ok := om.Get("x")
		if !ok || val != 15 {
			t.Errorf("Expected (15, true), got (%d, %v)", val, ok)
		}
	})

	t.Run("test delete", func(t *testing.T) {
		om.Delete("xyz") // no-op
		om.Delete("b")

		_, ok := om.Get("b")
		if ok {
			t.Error("Expected key to be deleted")
		}

		val, ok := om.Get("a")
		if !ok || val != 1 {
			t.Errorf("Expected 'a' to still exist with value 1, got (%d, %v)", val, ok)
		}

		val, ok = om.Get("x")
		if !ok || val != 15 {
			t.Errorf("Expected 'x' to still exist with value 15, got (%d, %v)", val, ok)
		}
	})
}

func TestOrderedMap_Prepend(t *testing.T) {
	om := gomaps.NewOrderedMap[string, int]()

	t.Run("test empty", func(t *testing.T) {
		err := om.Prepend("a", 1)
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}

		out, ok := om.Get("a")
		if !ok || out != 1 {
			t.Fatalf("expected (1 true), got (%d %t)", out, ok)
		}
	})

	t.Run("test non empty", func(t *testing.T) {
		om.Prepend("b", 2)
		om.Prepend("c", 3)
		om.Prepend("d", 4)

		expectedKeys := []string{"d", "c", "b", "a"}
		expectedValues := []int{4, 3, 2, 1}

		i := 0
		om.Range(func(key string, value int) bool {
			if expectedKeys[i] != key || expectedValues[i] != value {
				t.Fatalf("expected (%s %d), got (%s %d)", expectedKeys[i], expectedValues[i], key, value)
			}
			i++
			return true
		})
	})
}

func TestOrderedMap_Range(t *testing.T) {
	om := gomaps.NewOrderedMap[string, int]()

	t.Run("test empty", func(t *testing.T) {
		calls := 0
		om.Range(func(key string, value int) bool {
			calls++
			return true
		})

		if calls != 0 {
			t.Errorf("Expected 0 calls on empty map, got %d", calls)
		}
	})

	t.Run("test with values", func(t *testing.T) {
		om.Push("a", 1)
		om.Push("b", 2)
		om.Push("c", 3)

		expectedKeys := []string{"a", "b", "c"}
		expectedValues := []int{1, 2, 3}
		i := 0

		om.Range(func(key string, value int) bool {
			if key != expectedKeys[i] {
				t.Fatalf("Expected key %s at position %d, got %s", expectedKeys[i], i, key)
			}

			if value != expectedValues[i] {
				t.Errorf("Expected value %d at position %d, got %d", expectedValues[i], i, value)
			}

			i++

			return true
		})

		if i != len(expectedKeys) {
			t.Errorf("Expected %d iterations, got %d", len(expectedKeys), i)
		}

	})

	t.Run("test early stop", func(t *testing.T) {
		i := 0
		stopped := false

		om.Range(func(key string, value int) bool {
			if i == 1 {
				stopped = true
				return false
			}

			i++

			return true
		})

		if !stopped {
			t.Fatalf("Expected iteration to stop early but it didn't")
		}

		if i != 1 {
			t.Fatalf("Expected 1 call before stopping, got %d", i)
		}
	})
}

func TestOrderedMap_RangeUnsafe(t *testing.T) {
	om := gomaps.NewOrderedMap[int, string]()

	t.Run("test empty modify", func(t *testing.T) {
		calls := 0
		om.RangeUnsafe(func(key int, value string) bool {
			om.Set(key, value+" updated")
			calls++
			return true
		})

		if calls > 0 {
			t.Fatalf("Expected no calls, got %d", calls)
		}
	})

	t.Run("test modify with values", func(t *testing.T) {
		om.Set(1, "one")
		om.Set(2, "two")

		om.RangeUnsafe(func(k int, v string) bool {
			om.Set(k, v+" updated")
			return true
		})

		for _, key := range []int{1, 2} {
			val, ok := om.Get(key)
			if !ok {
				t.Fatalf("Expected key %d to exist", key)
				continue
			}

			if val != "one updated" && val != "two updated" {
				t.Fatalf("Expected value updated for key %d, got %s", key, val)
			}
		}
	})

	t.Run("test early stop", func(t *testing.T) {
		i := 0
		stopped := false

		om.RangeUnsafe(func(k int, v string) bool {
			if i == 1 {
				stopped = true
				return false
			}

			i++

			return true
		})

		if !stopped {
			t.Fatalf("Expected iteration to stop early but it didn't")
		}

		if i != 1 {
			t.Fatalf("Expected 1 call before stopping, got %d", i)
		}
	})
}

func TestOrderedMap_Reverse(t *testing.T) {
	om := gomaps.NewOrderedMap[int, string]()

	initialKeys := []int{1, 2, 3, 4}
	initialValues := []string{"one", "two", "three", "four"}
	reversedkeys := []int{4, 3, 2, 1}
	reversedValues := []string{"four", "three", "two", "one"}

	t.Run("initial reverse", func(t *testing.T) {
		om.Set(1, "one")
		om.Set(2, "two")
		om.Set(3, "three")
		om.Set(4, "four")

		om.Reverse()

		idx := 0
		om.Range(func(key int, value string) bool {
			if key != reversedkeys[idx] && value != reversedValues[idx] {
				t.Fatalf("expected %d: %s, got %d: %s", reversedkeys[idx], reversedValues[idx], key, value)
			}

			idx++

			return true
		})
	})

	t.Run("second reverse", func(t *testing.T) {
		om.Reverse()

		idx := 0
		om.Range(func(key int, value string) bool {
			if key != initialKeys[idx] && value != initialValues[idx] {
				t.Fatalf("expected %d: %s, got %d: %s", initialKeys[idx], initialValues[idx], key, value)
			}

			idx++

			return true
		})
	})

	t.Run("test empty", func(t *testing.T) {
		om := gomaps.NewOrderedMap[int, int]()

		defer func() {
			r := recover()
			if r != nil {
				t.Fatalf("expected success, got %v", r)
			}
		}()

		om.Reverse()
	})

	t.Run("test single element", func(t *testing.T) {
		om := gomaps.NewOrderedMap[string, int]()

		om.Set("only", 42)

		om.Reverse()

		value, ok := om.Get("only")
		if !ok || value != 42 {
			t.Fatalf("expected success, got ok=%t value=%d", ok, value)
		}
	})
}

func TestOrderedMap_Concat(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		om1 := gomaps.NewOrderedMap[string, int]()
		om2 := gomaps.NewOrderedMap[string, int]()
		out := om1.Concat(om2)

		i := 0
		out.Range(func(key string, value int) bool {
			i++
			return true
		})

		if i > 0 {
			t.Fatalf("Expected empty, got %d values", i)
		}
	})

	t.Run("test with values", func(t *testing.T) {
		om1 := gomaps.NewOrderedMap[string, int]()
		om1.Push("a", 1)
		om1.Push("b", 2)
		om1.Push("c", 3)

		om2 := gomaps.NewOrderedMap[string, int]()
		om2.Push("d", 4)
		om2.Push("e", 5)
		om2.Push("f", 6)

		expectedKeys := []string{"a", "b", "c", "d", "e", "f"}
		expectedValues := []int{1, 2, 3, 4, 5, 6}

		out := om1.Concat(om2)

		i := 0
		out.Range(func(key string, value int) bool {
			if key != expectedKeys[i] {
				t.Fatalf("Expected key %s at position %d, got %s", expectedKeys[i], i, key)
			}

			if value != expectedValues[i] {
				t.Errorf("Expected value %d at position %d, got %d", expectedValues[i], i, value)
			}

			i++

			return true
		})

		if i != len(expectedKeys) {
			t.Errorf("Expected %d iterations, got %d", len(expectedKeys), i)
		}
	})
}
