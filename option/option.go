package option

import (
	"encoding/json"
	"errors"
)

type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
}

func None[T any]() Option[T] {
	var zero T
	return Option[T]{value: zero, ok: false}
}

func (o Option[T]) IsSome() bool {
	return o.ok
}

func (o Option[T]) IsNone() bool {
	return !o.ok
}

func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("called Unwrap on a None value")
	}
	return o.value
}

func (o Option[T]) UnwrapOr(defaultVal T) T {
	if o.ok {
		return o.value
	}
	return defaultVal
}

func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.value
	}
	return f()
}

func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.ok {
		return Some(f(o.value))
	}
	return None[U]()
}

func (o Option[T]) And(other Option[T]) Option[T] {
	if o.ok {
		return other
	}
	return None[T]()
}

func (o Option[T]) Or(other Option[T]) Option[T] {
	if o.ok {
		return o
	}
	return other
}

func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.ok && predicate(o.value) {
		return o
	}
	return None[T]()
}

func (o Option[T]) Match(some func(T), none func()) {
	if o.ok {
		some(o.value)
	} else {
		none()
	}
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.ok {
		return json.Marshal(o.value)
	}
	return []byte("null"), nil
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.ok = false
		var zero T
		o.value = zero
		return nil
	}
	err := json.Unmarshal(data, &o.value)
	if err != nil {
		return err
	}
	o.ok = true
	return nil
}

func (o Option[T]) ToPointer() *T {
	if o.ok {
		return &o.value
	}
	return nil
}

func FromPointer[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

var ErrNone = errors.New("attempted to unwrap None Option")

func (o Option[T]) TryUnwrap() (T, error) {
	if o.ok {
		return o.value, nil
	}
	var zero T
	return zero, ErrNone
}

// Flatten converts Option[Option[T]] → Option[T]
func Flatten[T any](nested Option[Option[T]]) Option[T] {
	if nested.ok {
		return nested.value
	}
	return None[T]()
}
