package enumerable

type enumerableWhere[T any] struct {
	source    Enumerable[T]
	predicate func(T) (bool, error)
}

// Where creates an `Enumerable` from the given `Enumerable` and predicate. Items in the
// source `Enumerable` must return true when passed into the predicate in order to be yielded
// from the returned `Enumerable`.
func Where[T any](source Enumerable[T], predicate func(T) (bool, error)) Enumerable[T] {
	return &enumerableWhere[T]{
		source:    source,
		predicate: predicate,
	}
}

func (s *enumerableWhere[T]) Next() (bool, error) {
	for {
		hasNext, err := s.source.Next()
		if !hasNext || err != nil {
			return hasNext, err
		}

		value, err := s.source.Value()
		if err != nil {
			return false, err
		}

		if passes, err := s.predicate(value); passes || err != nil {
			return passes, err
		}
	}
}

func (s *enumerableWhere[T]) Value() (T, error) {
	return s.source.Value()
}

func (s *enumerableWhere[T]) Reset() {
	s.source.Reset()
}
