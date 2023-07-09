package app

const MODULE_NAME = "app"

type V1[T any] interface {
	Config() *T
}
