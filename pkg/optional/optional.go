package optional

import (
	"errors"
)

type Optional[T any] struct {
	v *T
}

func New[T any](v T) Optional[T] {
	return Optional[T]{
		v: &v,
	}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{
		v: nil,
	}
}

func (o Optional[T]) Present() bool {
	return o.v != nil
}

func (o Optional[T]) Get() (T, error) {
	if o.v == nil {
		var zero T
		return zero, errors.New("value not present")
	}
	return *o.v, nil
}

func (o Optional[T]) MustGet() T {
	v, err := o.Get()
	if err != nil {
		panic(err)
	}
	return v
}

func (o Optional[T]) IfPresent(fn func(T) T) Optional[T] {
	if o.Present() {
		return New(fn(o.MustGet()))
	}
	return o
}

func (o Optional[T]) OrElse(v T) Optional[T] {
	if o.Present() {
		return o
	}
	return New(v)
}

func Consumer[T any](fn func(T)) func(T) T {
	return func(v T) T {
		fn(v)
		return v
	}
}