package gomaps_test

import (
	"testing"

	"github.com/dimmerz92/gomaps"
)

func TestOrderedMap(t *testing.T) {
	var om *gomaps.OrderedMap[string, int]

	t.Run("test new ordered map", func(t *testing.T) {
		om = gomaps.NewOrderedMap[string, int]()
		if om == nil {
			t.Fatal("NewOrderedMap returned nil")
		}
	})

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
