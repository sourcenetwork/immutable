package immutable

import "encoding/json"

// Option represents an item that may or may not have a value.
type Option[T any] struct {
	// If HasValue is true, this Option contains a value, if
	// it is false it contains no value.
	hasValue bool

	// The Value of this Option. Should be ignored if HasValue is false.
	value T
}

// Some returns an `Option` of type `T` with the given value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		hasValue: true,
		value:    value,
	}
}

// Some returns an `Option` of type `T` with no value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// HasValue returns a boolean indicating whether or not this option contains a value. If
// it returns true, this Option contains a value, if it is false it contains no value.
func (o Option[T]) HasValue() bool {
	return o.hasValue
}

// Value returns the Value of this Option. Value returned is invalid HasValue() is false
// and should be ignored.
func (o Option[T]) Value() T {
	return o.value
}

// MarshalJSON implements the json.Marshaler interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.HasValue() {
		return json.Marshal(o.Value())
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *Option[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*o = None[T]()
		return nil
	}
	var value T
	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}
	*o = Some(value)
	return nil
}
