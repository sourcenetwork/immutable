package enumerable

type enumerableTake[T any] struct {
	source Enumerable[T]
	limit  uint64
	count  uint64
}

// Take creates an `Enumerable` from the given `Enumerable` and limit. The returned
// `Enumerable` will restrict the maximum number of items yielded to the given limit.
func Take[T any](source Enumerable[T], limit uint64) Enumerable[T] {
	return &enumerableTake[T]{
		source: source,
		limit:  limit,
	}
}

func (s *enumerableTake[T]) Next() (bool, error) {
	if s.count == s.limit {
		return false, nil
	}
	s.count += 1
	return s.source.Next()
}

func (s *enumerableTake[T]) Value() (T, error) {
	return s.source.Value()
}

func (s *enumerableTake[T]) Reset() {
	s.count = 0
	s.source.Reset()
}
