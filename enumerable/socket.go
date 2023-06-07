package enumerable

import "github.com/sourcenetwork/immutable"

// Socket is an extention of the enumerable interface allowing the source
// to be replaced after initial construction.
type Socket[T any] interface {
	Enumerable[T]
	// SetSource sets the source to this enumerable.
	//
	// This may be done after enumeration has begun.
	SetSource(Enumerable[T])
}

type socket[T any] struct {
	source immutable.Option[Enumerable[T]]
}

var _ Socket[any] = (*socket[any])(nil)

// NewSocket creates a new Socket enumerable with no initial source.
//
// The source may be set, and even swapped out, later during its lifetime.
// If enumeration begins before a source has been set it will behave as if empty.
// Reseting the Socket will reset the source if there is one, and then remove it
// as the source of this Socket.
func NewSocket[T any]() Socket[T] {
	return &socket[T]{
		source: immutable.None[Enumerable[T]](),
	}
}

// SetSource sets the source to this enumerable.
//
// This may be done after enumeration has begun.
func (s *socket[T]) SetSource(newSource Enumerable[T]) {
	s.source = immutable.Some(newSource)
}

func (s *socket[T]) Next() (bool, error) {
	if !s.source.HasValue() {
		return false, nil
	}

	return s.source.Value().Next()
}

func (s *socket[T]) Value() (T, error) {
	if !s.source.HasValue() {
		var v T
		return v, nil
	}

	return s.source.Value().Value()
}

func (s *socket[T]) Reset() {
	if s.source.HasValue() {
		s.source.Value().Reset()
	}
	s.source = immutable.None[Enumerable[T]]()
}
