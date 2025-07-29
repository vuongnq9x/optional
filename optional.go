// Package optional provides a generic Optional type that represents a value which may or may not be present.
// It helps avoid null pointer exceptions and provides a more expressive way to handle optional values.
package optional

import (
	"encoding/json"
	"fmt"
)

// Optional represents a value that may or may not be present
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional with a value
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// None creates an empty Optional
func None[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// FromPointer creates an Optional from a pointer.
// Returns None if the pointer is nil, or Some containing the pointed value if non-nil.
func FromPointer[T any](ptr *T) Optional[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// ToPointer converts the Optional to a pointer.
// Returns nil if the Optional is empty, or a pointer to the value if present.
func (o *Optional[T]) ToPointer() *T {
	if !o.present {
		return nil
	}
	return &o.value
}

// IsPresent returns true if the Optional contains a value
func (o *Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty returns true if the Optional is empty
func (o *Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value, panics if empty
func (o *Optional[T]) Get() T {
	if !o.present {
		panic("called Get() on empty Optional")
	}
	return o.value
}

// Equals checks if this Optional is equal to another Optional.
// Two Optionals are equal if they are both empty or contain equal values.
func (o *Optional[T]) Equals(other Optional[T]) bool {
	if o.present != other.present {
		return false
	}
	if !o.present {
		return true
	}
	return fmt.Sprintf("%v", o.value) == fmt.Sprintf("%v", other.value)
}

// Or returns this Optional if it has a value, otherwise returns the other Optional.
func (o *Optional[T]) Or(other Optional[T]) Optional[T] {
	if o.present {
		return *o
	}
	return other
}

// OrElsePanic returns the contained value if present, otherwise panics with the given message.
func (o *Optional[T]) OrElsePanic(message string) T {
	if !o.present {
		panic(message)
	}
	return o.value
}

// OrElse returns the value or a default value if empty
func (o *Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// OrElseGet returns the value or calls a supplier function if empty
func (o *Optional[T]) OrElseGet(supplier func() T) T {
	if o.present {
		return o.value
	}
	return supplier()
}

// IfPresent calls the consumer function if value is present
func (o *Optional[T]) IfPresent(consumer func(T)) {
	if o.present {
		consumer(o.value)
	}
}

// IfPresentOrElse executes the given consumer function if value is present,
// otherwise executes the runnable function.
func (o *Optional[T]) IfPresentOrElse(consumer func(T), runnable func()) {
	if o.present {
		consumer(o.value)
	} else {
		runnable()
	}
}

// Zip combines two Optionals into a single Optional using the provided combiner function.
// Returns None if either Optional is empty.
func Zip[T, U, R any](opt1 Optional[T], opt2 Optional[U], combiner func(T, U) R) Optional[R] {
	if opt1.present && opt2.present {
		return Some(combiner(opt1.value, opt2.value))
	}
	return None[R]()
}

// Map transforms the Optional value if present
func Map[T, U any](opt Optional[T], mapper func(T) U) Optional[U] {
	if opt.present {
		return Some(mapper(opt.value))
	}
	return None[U]()
}

// FlatMap transforms the Optional value to another Optional if present
func FlatMap[T, U any](opt Optional[T], mapper func(T) Optional[U]) Optional[U] {
	if opt.present {
		return mapper(opt.value)
	}
	return None[U]()
}

// Filter returns the Optional if the predicate is true, otherwise None
func (o *Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.present && predicate(o.value) {
		return *o
	}
	return None[T]()
}

// String returns string representation
func (o *Optional[T]) String() string {
	if o.present {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

// MarshalJSON implements json.Marshaler
func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.present {
		return json.Marshal(o.value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.present = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = value
	o.present = true
	return nil
}
