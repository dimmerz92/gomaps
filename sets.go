package gomaps

type Set[T comparable] map[T]struct{}

func ToSet[T comparable](slice []T) Set[T] {
	set := make(Set[T])

	for _, value := range slice {
		set[value] = struct{}{}
	}

	return set
}

// Union returns the union of any number of sets.
func (s *Set[T]) Union(sets ...Set[T]) Set[T] {
	result := make(Set[T])

	for _, set := range append(sets, *s) {
		for k := range set {
			result[k] = struct{}{}
		}
	}

	return result
}

// Intersect returns the intersection of any number of sets.
func (s *Set[T]) Intersect(sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return Set[T]{}
	}

	minSetIdx := 0
	minSize := len(sets[0])
	sets = append(sets, *s)

	for i, set := range sets {
		if len(set) < minSize {
			minSetIdx = i
			minSize = len(set)
		}
	}

	result := make(Set[T])
	for value := range sets[minSetIdx] {
		result[value] = struct{}{}
	}

	for i, set := range sets {
		if i == minSetIdx {
			continue
		}

		for value := range result {
			_, exists := set[value]
			if !exists {
				delete(result, value)
			}
		}
	}

	return result
}

func (s *Set[T]) Difference(set Set[T]) Set[T] {
	result := *s

	for key := range result {
		_, exists := set[key]
		if exists {
			delete(result, key)
		}
	}

	return result
}
