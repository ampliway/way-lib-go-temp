package v1

import (
	"errors"
	"fmt"

	"github.com/ampliway/way-lib-go/app"
	"github.com/ampliway/way-lib-go/config"
	configV1 "github.com/ampliway/way-lib-go/config/v1"
)

var (
	_                app.V1[any] = (*App[any])(nil)
	errSubModuleInit             = errors.New("sub-module failed on init")
)

type App[T any] struct {
	config config.V1[T]
}

func New[T any]() (*App[T], error) {
	config, err := configV1.New[T]()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errSubModuleInit, err)
	}

	return &App[T]{
		config: config,
	}, nil
}

func (a *App[T]) Config() *T {
	return a.config.Get()
}
