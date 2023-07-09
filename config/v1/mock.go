package v1

import (
	"github.com/ampliway/way-lib-go/config"
)

var _ config.V1[any] = (*Mock[any])(nil)

type Mock[T any] struct {
	value *T
}

func NewMock[T any](value T) (*Mock[T], error) {
	return &Mock[T]{
		value: &value,
	}, nil
}

func (e *Mock[T]) Get() *T {
	return e.value
}
