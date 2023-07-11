package app

import "github.com/ampliway/way-lib-go/msg"

const MODULE_NAME = "app"

type V1[T any] interface {
	Config() *T
	Msg() msg.ProducerV1
}
