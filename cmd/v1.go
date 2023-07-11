package cmd

import "github.com/ampliway/way-lib-go/app"

const MODULE_NAME = "cmd"

type V1[T any] interface {
	Add(config *Config[T]) error
	Run(arguments ...string)
}

type Config[T any] struct {
	Name        string
	Description string
	Execute     func(app app.V1[T]) error
	Args        []string
}
