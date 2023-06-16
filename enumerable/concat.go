package enumerable

// Concatenation is an extention of the enumerable interface allowing new sources
// to be added after initial construction.
type Concatenation[T any] interface {
	Enumerable[T]
	// Append appends a new source to this concatenation.
	//
	// This may be done after enumeration has begun.
	Append(Enumerable[T])
}

type enumerableConcat[T any] struct {
	sources            []Enumerable[T]
	currentSourceIndex int
}

// Concat takes zero to many source `Ènumerable`s and stacks them on top
// of each other, resulting in one enumerable that will iterate through all
// the values in all of the given sources.
//
// New sources may be added after iteration has begun.
func Concat[T any](sources ...Enumerable[T]) Concatenation[T] {
	return &enumerableConcat[T]{
		sources:            sources,
		currentSourceIndex: 0,
	}
}

// Append appends a new source to this concatenation.
//
// This may be done after enumeration has begun.
func (s *enumerableConcat[T]) Append(newSource Enumerable[T]) {
	s.sources = append(s.sources, newSource)
}

func (s *enumerableConcat[T]) Next() (bool, error) {
	startSourceIndex := s.currentSourceIndex
	hasLooped := false

	for {
		if s.currentSourceIndex >= len(s.sources) {
			if len(s.sources) < 1 || hasLooped {
				return false, nil
			}
			s.currentSourceIndex = 0
			hasLooped = true
		}

		currentSource := s.sources[s.currentSourceIndex]
		hasValue, err := currentSource.Next()
		if err != nil {
			return false, err
		}
		if hasValue {
			return true, nil
		}

		s.currentSourceIndex += 1

		if s.currentSourceIndex == startSourceIndex {
			return false, nil
		}
	}
}

func (s *enumerableConcat[T]) Value() (T, error) {
	return s.sources[s.currentSourceIndex].Value()
}

func (s *enumerableConcat[T]) Reset() {
	s.currentSourceIndex = 0
	for _, source := range s.sources {
		source.Reset()
	}
}
