package result

import (
	"encoding/json"
)

type Result[T any] struct {
	value T
	err   error
	ok    bool
}

func Ok[T any](v T) Result[T] {
	return Result[T]{value: v, err: nil, ok: true}
}

func Err[T any](e error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: e, ok: false}
}

func (r Result[T]) IsOk() bool  { return r.ok }
func (r Result[T]) IsErr() bool { return !r.ok }

func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic("called Unwrap on an Err value")
	}
	return r.value
}

func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.ok {
		return r.value
	}
	return defaultVal
}

func (r Result[T]) UnwrapErr() error {
	if r.ok {
		panic("called UnwrapErr on an Ok value")
	}
	return r.err
}

func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.ok {
		return Ok[U](f(r.value))
	}
	return Err[U](r.err)
}

func MapErr[T any](r Result[T], f func(error) error) Result[T] {
	if r.ok {
		return Ok[T](r.value)
	}
	return Err[T](f(r.err))
}

func (r Result[T]) Match(ok func(T), err func(error)) {
	if r.ok {
		ok(r.value)
	} else {
		err(r.err)
	}
}

func (r Result[T]) AndThen(f func(T) Result[T]) Result[T] {
	if r.ok {
		return f(r.value)
	}
	return r
}

func FlatMap[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.ok {
		return f(r.value)
	}
	return Err[U](r.err)
}

func Flatten[T any](nested Result[Result[T]]) Result[T] {
	if nested.IsErr() {
		return Err[T](nested.UnwrapErr())
	}
	return nested.Unwrap()
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.ok {
		return json.Marshal(map[string]any{
			"ok":    true,
			"value": r.value,
		})
	}
	return json.Marshal(map[string]any{
		"ok":  false,
		"err": r.err.Error(),
	})
}
