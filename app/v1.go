package app

import (
	"github.com/ampliway/way-lib-go/cache"
	"github.com/ampliway/way-lib-go/msg"
	"github.com/ampliway/way-lib-go/storage"
)

const MODULE_NAME = "app"

type V1[T any] interface {
	Config() *T
	Msg() msg.MsgV1
	Storage() storage.V1
	Cache() cache.V1
	ID() string
}
