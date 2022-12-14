package enumerable

import "sort"

type enumerableSort[T any] struct {
	source   Enumerable[T]
	less     func(T, T) bool
	capacity int
	result   Enumerable[T]
}

// Sort creates an `Enumerable` from the given `Enumerable`, using the given
// less function to determine as to whether an item is less than the other in
// in terms of order.
//
// The returned `Enumerable` will enumerate the entire source
// enumerable on the first `Next` call, but will not enumerate it again unless
// reset.
func Sort[T any](source Enumerable[T], less func(T, T) bool, capacity int) Enumerable[T] {
	return &enumerableSort[T]{
		source:   source,
		less:     less,
		capacity: capacity,
	}
}

func (s *enumerableSort[T]) Next() (bool, error) {
	if s.result == nil {
		result := make([]T, 0, s.capacity)
		var searchErr error

		// Declaring an anonymous function costs, so we do it here outside of the loop
		// even though it is slightly less intuitive
		f := func(i int) bool {
			var val T
			val, searchErr = s.source.Value()
			return !s.less(result[i], val)
		}

		for i := 0; i <= s.capacity; i++ {
			hasNext, err := s.source.Next()
			if err != nil {
				return false, err
			}
			if !hasNext {
				break
			}

			previousLength := len(result)
			indexOfFirstGreaterValue := sort.Search(previousLength, f)
			if searchErr != nil {
				// This is quite ugly, but sort.Search does not allow for anything else
				return false, searchErr
			}

			value, err := s.source.Value()
			result = append(result, value)
			if indexOfFirstGreaterValue == previousLength {
				// Value is the greatest, and belongs at the end
				continue
			}
			// Shift all items to the right of the first element of greater value by
			// one place.  This call should not allocate.
			copy(result[indexOfFirstGreaterValue+1:], result[indexOfFirstGreaterValue:])
			result[indexOfFirstGreaterValue] = value
		}

		// Use the enumerableSlice for convienience
		s.result = New(result)
	}

	return s.result.Next()
}

func (s *enumerableSort[T]) Value() (T, error) {
	return s.result.Value()
}

func (s *enumerableSort[T]) Reset() {
	// s.result should be cleared, not reset, as Reset should
	// enable the re-enumeration of the entire enumeration chain,
	// not just the last step.
	s.result = nil
	s.source.Reset()
}
