package gomaps_test

import (
	"reflect"
	"testing"

	"github.com/dimmerz92/gomaps"
)

func TestSet_ToSet(t *testing.T) {
	t.Run("test strings", func(t *testing.T) {
		in := []string{"1", "2", "3"}
		expected := gomaps.Set[string]{"1": {}, "2": {}, "3": {}}

		out := gomaps.ToSet(in)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})
}

func TestSet_Union(t *testing.T) {
	t.Run("test 3x empty sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{}
		s2 := gomaps.Set[int]{}
		s3 := gomaps.Set[int]{}

		expected := gomaps.Set[int]{}

		out := s1.Union(s2, s3)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})

	t.Run("test 2x empty, 1x full sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{}
		s2 := gomaps.Set[int]{1: {}, 2: {}, 3: {}}
		s3 := gomaps.Set[int]{}

		expected := gomaps.Set[int]{1: {}, 2: {}, 3: {}}

		out := s1.Union(s2, s3)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})

	t.Run("test 3x full sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{1: {}, 2: {}, 3: {}}
		s2 := gomaps.Set[int]{2: {}, 3: {}, 4: {}}
		s3 := gomaps.Set[int]{3: {}, 4: {}, 5: {}}

		expected := gomaps.Set[int]{1: {}, 2: {}, 3: {}, 4: {}, 5: {}}

		out := s1.Union(s2, s3)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})
}

func TestSet_Intersect(t *testing.T) {
	t.Run("test 3x empty sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{}
		s2 := gomaps.Set[int]{}
		s3 := gomaps.Set[int]{}

		expected := gomaps.Set[int]{}

		out := s1.Intersect(s2, s3)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})

	t.Run("test 2x empty, 1x full sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{}
		s2 := gomaps.Set[int]{1: {}, 2: {}, 3: {}}
		s3 := gomaps.Set[int]{}

		expected := gomaps.Set[int]{}

		out := s1.Intersect(s2, s3)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})

	t.Run("test 3x full sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{1: {}, 2: {}, 3: {}}
		s2 := gomaps.Set[int]{2: {}, 3: {}, 4: {}}
		s3 := gomaps.Set[int]{3: {}, 4: {}, 5: {}}

		expected := gomaps.Set[int]{3: {}}

		out := s1.Intersect(s2, s3)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})
}

func TestSet_Difference(t *testing.T) {
	t.Run("test 2x empty sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{}
		s2 := gomaps.Set[int]{}

		expected := gomaps.Set[int]{}

		out := s1.Difference(s2)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})

	t.Run("test 1x empty, 1x full sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{}
		s2 := gomaps.Set[int]{1: {}, 2: {}, 3: {}}

		expected := gomaps.Set[int]{}

		out := s1.Intersect(s2)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})

	t.Run("test 2x full sets", func(t *testing.T) {
		s1 := gomaps.Set[int]{1: {}, 2: {}, 3: {}}
		s2 := gomaps.Set[int]{2: {}, 3: {}, 4: {}}

		expected := gomaps.Set[int]{1: {}}

		out := s1.Difference(s2)

		if !reflect.DeepEqual(out, expected) {
			t.Fatalf("expected %#v, got %#v", expected, out)
		}
	})
}
