package config

const MODULE_NAME = "config"

type V1[T any] interface {
	Get() *T
}
