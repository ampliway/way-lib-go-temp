package app

import "github.com/ampliway/way-lib-go/config"

const MODULE_NAME = "app"

type V1[T any] interface {
	Config() config.V1[T]
}
