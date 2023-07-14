package v1

import (
	"errors"
	"fmt"

	"github.com/ampliway/way-lib-go/app"
	"github.com/ampliway/way-lib-go/config"
	configV1 "github.com/ampliway/way-lib-go/config/v1"
	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/msg"
	msgV1 "github.com/ampliway/way-lib-go/msg/v1"
	"github.com/ampliway/way-lib-go/storage"
	storageV1 "github.com/ampliway/way-lib-go/storage/v1"
)

var (
	_                app.V1[any] = (*App[any])(nil)
	errSubModuleInit             = errors.New("sub-module failed on init")
)

type App[T any] struct {
	config  config.V1[T]
	msg     msg.ProducerV1
	storage storage.V1
}

func New[T any]() (*App[T], error) {
	cfg, err := configV1.New[T]()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", errSubModuleInit, err, config.MODULE_NAME)
	}

	natsConfig, err := configV1.New[msgV1.Config]()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", errSubModuleInit, err, msg.MODULE_NAME)
	}

	m, err := msgV1.New(natsConfig.Get(), id.New())
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", errSubModuleInit, err, msg.MODULE_NAME)
	}

	storageConfig, err := configV1.New[storageV1.Config]()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", errSubModuleInit, err, msg.MODULE_NAME)
	}

	s, err := storageV1.New(storageConfig.Get(), id.New())
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s", errSubModuleInit, err, storage.MODULE_NAME)
	}

	return &App[T]{
		config:  cfg,
		msg:     m,
		storage: s,
	}, nil
}

func (a *App[T]) Config() *T {
	return a.config.Get()
}

func (a *App[T]) Msg() msg.ProducerV1 {
	return a.msg
}

func (a *App[T]) Storage() storage.V1 {
	return a.storage
}
