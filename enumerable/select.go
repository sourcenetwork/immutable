package enumerable

type enumerableSelect[TSource any, TResult any] struct {
	source       Enumerable[TSource]
	selector     func(TSource) (TResult, error)
	currentValue TResult
}

// Select creates a new `Enumerable` that iterates through each item
// yielded by the given source and then yields the value returned by
// the given selector.
func Select[TSource any, TResult any](
	source Enumerable[TSource],
	selector func(TSource) (TResult, error),
) Enumerable[TResult] {
	return &enumerableSelect[TSource, TResult]{
		source:   source,
		selector: selector,
	}
}

func (s *enumerableSelect[TSource, TResult]) Next() (bool, error) {
	hasNext, err := s.source.Next()
	if !hasNext || err != nil {
		return hasNext, err
	}

	value, err := s.source.Value()
	if err != nil {
		return false, nil
	}

	// We do this here to keep the work (and errors) in the `Next` call
	result, err := s.selector(value)
	if err != nil {
		return false, nil
	}

	s.currentValue = result
	return true, nil
}

func (s *enumerableSelect[TSource, TResult]) Value() (TResult, error) {
	return s.currentValue, nil
}

func (s *enumerableSelect[TSource, TResult]) Reset() {
	s.source.Reset()
}
